package spotify

import (
	"context"
	"fmt"

	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/pkg/utils"
)

func (c *client) playlistSync(ctx context.Context, user model.User) error {
	playlistsDB, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	playlistsSpotify, err := c.playlistGetAll(ctx, user)
	if err != nil {
		return err
	}

	toCreate := make([]model.Playlist, 0)
	toUpdate := make([]model.Playlist, 0)
	toDelete := make([]model.Playlist, 0)

	// Find the playlists that need to be created or updated
	for i := range playlistsSpotify {
		playlistDB, ok := utils.SliceFind(playlistsDB, func(p *model.Playlist) bool { return p.Equal(playlistsSpotify[i].model) })
		if !ok {
			// Playlist doesn't exist yet
			// Create it
			toCreate = append(toCreate, playlistsSpotify[i].model)
			if err := c.playlistSaveCover(&playlistsSpotify[i].model, nil, playlistsSpotify[i].Images); err != nil {
				return err
			}

			continue
		}

		// Regardless if any of the other data changed, let's update the cover if we can
		if err := c.playlistSaveCover(&playlistsSpotify[i].model, *playlistDB, playlistsSpotify[i].Images); err != nil {
			return err
		}

		// Playlist already exist
		// But is it still completely the same?
		if !(*playlistDB).EqualEntry(playlistsSpotify[i].model) {
			// Not completely the same anymore
			// Update it
			toUpdate = append(toUpdate, playlistsSpotify[i].model)
		}
	}

	for i := range toCreate {
		if err := c.playlistUserCheck(ctx, toCreate[i].OwnerUID); err != nil {
			return err
		}
		if err := c.playlist.Create(ctx, &toCreate[i]); err != nil {
			return err
		}
	}

	for i := range toUpdate {
		if err := c.playlistUserCheck(ctx, toUpdate[i].OwnerUID); err != nil {
			return err
		}
		if err := c.playlist.Update(ctx, toUpdate[i]); err != nil {
			return err
		}
	}

	// New and updated entries are now in the database
	// Let's bring our local copy up to date
	playlistsDB, err = c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	// Find the playlists that need to be deleted
	for _, playlistDB := range playlistsDB {
		_, ok := utils.SliceFind(playlistsSpotify, func(p playlist) bool { return p.model.SpotifyID == playlistDB.SpotifyID })
		if !ok {
			// Playlist no longer exists in the user's account
			// So delete it
			toDelete = append(toDelete, *playlistDB)
		}
	}

	for i := range toDelete {
		if err := c.playlist.Delete(ctx, toDelete[i].ID); err != nil {
			return err
		}
	}

	return nil
}

// playlistTrackSync brings the local database up to date with the songs for each playlist
// nolint:gocognit // It's fine
func (c *client) playlistTrackSync(ctx context.Context, user model.User) error {
	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}
	if playlists == nil {
		return nil
	}

	for _, playlist := range playlists {
		tracksDB, err := c.track.GetByPlaylist(ctx, playlist.ID)
		if err != nil {
			return err
		}

		tracksSpotify, err := c.playlistGetTracksAll(ctx, user, *playlist)
		if err != nil {
			return err
		}

		toCreate := make([]model.Track, 0)
		toDelete := make([]model.Track, 0)

		for _, trackSpotify := range tracksSpotify {
			if _, ok := utils.SliceFind(tracksDB, func(t *model.Track) bool { return t.Equal(trackSpotify) }); !ok {
				toCreate = append(toCreate, trackSpotify)
			}
		}

		for _, trackDB := range tracksDB {
			if _, ok := utils.SliceFind(tracksSpotify, func(t model.Track) bool { return t.Equal(*trackDB) }); !ok {
				toDelete = append(toDelete, *trackDB)
			}
		}

		// Do the db operations
		for _, track := range toCreate {
			trackDB, err := c.track.GetBySpotify(ctx, track.SpotifyID)
			if err != nil {
				return err
			}

			if trackDB == nil {
				// We don't have the track yet
				if err := c.track.Create(ctx, &track); err != nil {
					return err
				}
				trackDB = &model.Track{
					ID: track.ID,
				}
			} else if !trackDB.Equal(track) {
				// Track is not up to date
				if err := c.track.UpdateBySpotify(ctx, track); err != nil {
					return err
				}
			}

			if err := c.playlist.CreateTrack(ctx, &model.PlaylistTrack{
				PlaylistID: playlist.ID,
				TrackID:    trackDB.ID,
			}); err != nil {
				return err
			}
		}

		for _, track := range toDelete {
			if err := c.playlist.DeleteTrackByPlaylistTrack(ctx, model.PlaylistTrack{
				PlaylistID: playlist.ID,
				TrackID:    track.ID,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// userSync updates the information for every relevant user (for the given user)
func (c *client) userSync(ctx context.Context, user model.User) error {
	// Get all relevant users
	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}
	if playlists == nil {
		return nil
	}

	usersDB := utils.SliceMap(playlists, func(p *model.Playlist) model.User { return p.Owner })
	usersDB = utils.SliceUnique(usersDB)

	// Get all spotify users
	usersSpotify := make([]model.User, 0, len(usersDB))
	for _, userDB := range usersDB {
		newUser, err := c.userGet(ctx, user, userDB)
		if err != nil {
			return err
		}

		usersSpotify = append(usersSpotify, newUser)
	}

	toUpdate := make([]model.User, 0)

	for _, userSpotify := range usersSpotify {
		if _, ok := utils.SliceFind(usersDB, func(u model.User) bool { return u.Equal(userSpotify) }); !ok {
			toUpdate = append(toUpdate, userSpotify)
		}
	}

	for _, user := range toUpdate {
		if err := c.user.Update(ctx, user); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) linkSync(ctx context.Context, user model.User) error {
	directories, err := c.directory.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	links, err := c.link.GetAllByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, link := range links {
		var sources []model.Playlist
		var targets []model.Playlist

		switch {
		case link.SourceDirectoryID != 0:
			directory, ok := utils.SliceFind(directories, func(d *model.Directory) bool { return d.ID == link.SourceDirectoryID })
			if !ok {
				return fmt.Errorf("database foreign key reference error (source directory) for link %+v", *link)
			}
			sources = (*directory).Playlists

		case link.SourcePlaylistID != 0:
			playlist, ok := utils.SliceFind(playlists, func(p *model.Playlist) bool { return p.ID == link.SourcePlaylistID })
			if !ok {
				return fmt.Errorf("database foreign key reference error (source playlist) for link %+v", *link)
			}
			sources = []model.Playlist{**playlist}

		default:
			return fmt.Errorf("database foreign key reference error (source) for link %+v", *link)
		}

		switch {
		case link.TargetDirectoryID != 0:
			directory, ok := utils.SliceFind(directories, func(d *model.Directory) bool { return d.ID == link.TargetDirectoryID })
			if !ok {
				return fmt.Errorf("database foreign key reference error (target directory) for link %+v", *link)
			}
			targets = (*directory).Playlists

		case link.TargetPlaylistID != 0:
			playlist, ok := utils.SliceFind(playlists, func(p *model.Playlist) bool { return p.ID == link.TargetPlaylistID })
			if !ok {
				return fmt.Errorf("database foreign key reference error (target playlist) for link %+v", *link)
			}
			targets = []model.Playlist{**playlist}

		default:
			return fmt.Errorf("database foreign key reference error (target) for link %+v", *link)
		}

		for i := range sources {
			for j := range targets {
				if err := c.linkSyncOne(ctx, user, sources[i], targets[j]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *client) linkSyncOne(ctx context.Context, user model.User, source, target model.Playlist) error {
	if source.Equal(target) {
		return nil
	}

	tracksSource, err := c.track.GetByPlaylist(ctx, source.ID)
	if err != nil {
		return err
	}

	tracksTarget, err := c.track.GetByPlaylist(ctx, target.ID)
	if err != nil {
		return err
	}

	toAdd := make([]model.Track, 0)

	for _, trackSource := range tracksSource {
		if _, ok := utils.SliceFind(tracksTarget, func(t *model.Track) bool { return t.Equal(*trackSource) }); !ok {
			toAdd = append(toAdd, *trackSource)
		}
	}

	if err := c.playlistAddTracksAll(ctx, user, target, toAdd); err != nil {
		return err
	}

	return nil
}
