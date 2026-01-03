package spotifysync

import (
	"context"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/pkg/utils"
)

// albumSync will syncronize the user's saved albums
func (c *client) albumSync(ctx context.Context, user model.User) error {
	albumsDB, err := c.album.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	albumsSpotifyAPI, err := spotifyapi.C.AlbumGetUser(ctx, user)
	if err != nil {
		return err
	}
	albumsSpotify := utils.SliceMap(albumsSpotifyAPI, func(a spotifyapi.Album) model.Album { return a.ToModel() })

	return syncUserData(syncUserDataStruct[model.Album]{
		DB:     utils.SliceDereference(albumsDB),
		API:    albumsSpotify,
		Equal:  func(a1, a2 model.Album) bool { return a1.Equal(a2) },
		Get:    func(a model.Album) (*model.Album, error) { return c.album.GetBySpotify(ctx, a.SpotifyID) },
		Create: func(a *model.Album) error { return c.album.Create(ctx, a) },
		CreateUserLink: func(a model.Album) error {
			return c.album.CreateUser(ctx, &model.AlbumUser{AlbumID: a.ID, UserID: user.ID})
		},
		DeleteUserLink: func(a model.Album) error {
			return c.album.DeleteUserByUserAlbum(ctx, model.AlbumUser{AlbumID: a.ID, UserID: user.ID})
		},
	})
}

// albumUpdate updates local album instances to match the spotify data.
// It updates all albums, regardless of the user given.
// However the given user's access token is used.
func (c *client) albumUpdate(ctx context.Context, user model.User) error {
	albumsDB, err := c.album.GetAll(ctx)
	if err != nil {
		return err
	}

	filtered := filterSpotify(filterSpotifyStruct[*model.Album]{
		Items:     albumsDB,
		Frequency: 24,
		SpotifyID: func(a *model.Album) string { return a.SpotifyID },
		UpdatedAt: func(a *model.Album) time.Time { return a.UpdatedAt },
	})
	if len(filtered) == 0 {
		return nil
	}

	albumsSpotifyAPI, err := spotifyapi.C.AlbumGetAll(ctx, user, filtered)
	if err != nil {
		return err
	}
	albumsSpotify := utils.SliceMap(albumsSpotifyAPI, func(a spotifyapi.Album) model.Album { return a.ToModel() })

	for i := range albumsSpotify {
		albumDB, ok := utils.SliceFind(albumsDB, func(a *model.Album) bool { return a.Equal(albumsSpotify[i]) })
		if !ok {
			// Album not found
			continue
		}

		albumsSpotify[i].ID = (*albumDB).ID

		// Bring the album data up to date
		a := albumsSpotify[i]
		if (*albumDB).EqualEntry(a) {
			a = model.Album{ID: a.ID} // Do an empty update to refresh updated_at
		}
		if err := c.album.Update(ctx, a); err != nil {
			return err
		}

		// Bring the album artists up to date
		artistsDB, err := c.artist.GetByAlbum(ctx, (*albumDB).ID)
		if err != nil {
			return err
		}

		artistsSpotify := utils.SliceMap(albumsSpotifyAPI[i].Artists, func(a spotifyapi.Artist) model.Artist { return a.ToModel() })

		if err := syncUserData(syncUserDataStruct[model.Artist]{
			DB:     utils.SliceDereference(artistsDB),
			API:    artistsSpotify,
			Equal:  func(a1, a2 model.Artist) bool { return a1.Equal(a2) },
			Get:    func(a model.Artist) (*model.Artist, error) { return c.artist.GetBySpotify(ctx, a.SpotifyID) },
			Create: func(a *model.Artist) error { return c.artist.Create(ctx, a) },
			CreateUserLink: func(a model.Artist) error {
				return c.album.CreateArtist(ctx, &model.AlbumArtist{AlbumID: (*albumDB).ID, ArtistID: a.ID})
			},
			DeleteUserLink: func(a model.Artist) error {
				return c.album.DeleteArtistByArtistAlbum(ctx, model.AlbumArtist{AlbumID: (*albumDB).ID, ArtistID: a.ID})
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) albumCoverSync(ctx context.Context, user model.User) error {
	albums, err := c.album.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	return c.syncCover(ctx, utils.SliceMap(albums, func(a *model.Album) syncCoverStruct {
		return syncCoverStruct{
			CoverURL: a.CoverURL,
			CoverID:  a.CoverID,
			Update: func(newID string) error {
				a.CoverID = newID
				return c.album.Update(ctx, *a)
			},
		}
	}))
}
