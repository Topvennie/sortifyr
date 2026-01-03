package spotifysync

import (
	"context"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/pkg/utils"
)

func (c *client) playlistSync(ctx context.Context, user model.User) error {
	playlistsDB, err := c.playlist.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	playlistsSpotifyAPI, err := spotifyapi.C.PlaylistGetUser(ctx, user)
	if err != nil {
		return err
	}
	playlistsSpotify := make([]model.Playlist, 0, len(playlistsSpotifyAPI))
	for i := range playlistsSpotifyAPI {
		ownerDB, err := c.user.GetByUID(ctx, playlistsSpotifyAPI[i].Owner.UID)
		if err != nil {
			return err
		}
		if ownerDB == nil {
			ownerDB = &model.User{UID: playlistsSpotifyAPI[i].Owner.UID, DisplayName: playlistsSpotifyAPI[i].Owner.DisplayName}
			if err := c.user.Create(ctx, ownerDB); err != nil {
				return err
			}
		}

		playlistSpotify := playlistsSpotifyAPI[i].ToModel()
		playlistSpotify.OwnerID = ownerDB.ID
		playlistSpotify.SnapshotID = "" // Make sure it gets fully updated

		playlistsSpotify = append(playlistsSpotify, playlistSpotify)
	}

	return syncUserData(syncUserDataStruct[model.Playlist]{
		DB:     utils.SliceDereference(playlistsDB),
		API:    playlistsSpotify,
		Equal:  func(p1, p2 model.Playlist) bool { return p1.Equal(p2) },
		Get:    func(p model.Playlist) (*model.Playlist, error) { return c.playlist.GetBySpotify(ctx, p.SpotifyID) },
		Create: func(p *model.Playlist) error { return c.playlist.Create(ctx, p) },
		CreateUserLink: func(p model.Playlist) error {
			return c.playlist.CreateUser(ctx, &model.PlaylistUser{PlaylistID: p.ID, UserID: user.ID})
		},
		DeleteUserLink: func(p model.Playlist) error {
			return c.playlist.DeleteUserByUserPlaylist(ctx, model.PlaylistUser{PlaylistID: p.ID, UserID: user.ID})
		},
	})
}

// playlistUpdate updates local playlist instances to match the spotify data
func (c *client) playlistUpdate(ctx context.Context, user model.User) error {
	playlistsDB, err := c.playlist.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	playlistsSpotifyAPI, err := spotifyapi.C.PlaylistGetUser(ctx, user)
	if err != nil {
		return err
	}
	playlistsSpotify := make([]model.Playlist, 0, len(playlistsSpotifyAPI))
	for i := range playlistsSpotifyAPI {
		ownerDB, err := c.user.GetByUID(ctx, playlistsSpotifyAPI[i].Owner.UID)
		if err != nil {
			return err
		}
		if ownerDB == nil {
			ownerDB = &model.User{UID: playlistsSpotifyAPI[i].Owner.UID, DisplayName: playlistsSpotifyAPI[i].Owner.DisplayName}
			if err := c.user.Create(ctx, ownerDB); err != nil {
				return err
			}
		}

		playlistSpotify := playlistsSpotifyAPI[i].ToModel()
		playlistSpotify.OwnerID = ownerDB.ID

		playlistsSpotify = append(playlistsSpotify, playlistSpotify)
	}

	for i := range playlistsSpotify {
		playlistDB, ok := utils.SliceFind(playlistsDB, func(p *model.Playlist) bool { return p.Equal(playlistsSpotify[i]) })
		if !ok {
			// Playlist not found
			continue
		}

		if (*playlistDB).SnapshotID == playlistsSpotify[i].SnapshotID {
			// Same snapshot id
			// Meaning they are still equal
			continue
		}

		playlistsSpotify[i].ID = (*playlistDB).ID

		// bring the playlist up to date
		if !(*playlistDB).EqualEntry(playlistsSpotify[i]) {
			if err := c.playlist.Update(ctx, playlistsSpotify[i]); err != nil {
				return err
			}
		}

		// Bring the playlist tracks up to date
		tracksDB, err := c.track.GetByPlaylist(ctx, (*playlistDB).ID)
		if err != nil {
			return err
		}

		tracksSpotifyAPI, err := spotifyapi.C.PlaylistGetTrackAll(ctx, user, playlistsSpotify[i].SpotifyID)
		if err != nil {
			return err
		}
		tracksSpotify := utils.SliceMap(tracksSpotifyAPI, func(t spotifyapi.Track) model.Track { return t.ToModel() })

		if err := syncUserData(syncUserDataStruct[model.Track]{
			DB:  utils.SliceDereference(tracksDB),
			API: tracksSpotify,
			Equal: func(t1, t2 model.Track) bool {
				// Unavailable tracks have no spotify id
				// So we need another method to check for equality if that is the case
				if t1.SpotifyID != "" || t2.SpotifyID != "" {
					return t1.Equal(t2)
				}

				return t1.Name == t2.Name
			},
			Get: func(t model.Track) (*model.Track, error) {
				// Again we need to support tracks without spotify id
				if t.SpotifyID != "" {
					return c.track.GetBySpotify(ctx, t.SpotifyID)
				}

				tracks, err := c.track.GetByName(ctx, t.Name)
				if err != nil {
					return nil, err
				}

				// Only look for other tracks without spotify id (unavailable)
				if tt, ok := utils.SliceFind(tracks, func(t *model.Track) bool { return t.SpotifyID == "" }); ok {
					return *tt, nil
				}

				return nil, nil
			},
			Create: func(t *model.Track) error { return c.track.Create(ctx, t) },
			CreateUserLink: func(t model.Track) error {
				return c.playlist.CreateTrack(ctx, &model.PlaylistTrack{PlaylistID: (*playlistDB).ID, TrackID: t.ID})
			},
			DeleteUserLink: func(t model.Track) error {
				return c.playlist.DeleteTrackByPlaylistTrack(ctx, model.PlaylistTrack{PlaylistID: (*playlistDB).ID, TrackID: t.ID})
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) playlistCoverSync(ctx context.Context, user model.User) error {
	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	return c.syncCover(ctx, utils.SliceMap(playlists, func(p *model.Playlist) syncCoverStruct {
		return syncCoverStruct{
			CoverURL: p.CoverURL,
			CoverID:  p.CoverID,
			Update: func(newID string) error {
				p.CoverID = newID
				return c.playlist.Update(ctx, *p)
			},
		}
	}))
}
