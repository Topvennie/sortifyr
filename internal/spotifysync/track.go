package spotifysync

import (
	"context"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/pkg/utils"
)

// trackUpdate updates local track instances to match the spotify data.
// It updates all tracks, regardless of the user given.
// However the given user's access token is used.
func (c *client) trackUpdate(ctx context.Context, user model.User) error {
	tracksDB, err := c.track.GetAll(ctx)
	if err != nil {
		return err
	}

	filtered := filterSpotify(filterSpotifyStruct[*model.Track]{
		Items:     tracksDB,
		Frequency: 48,
		SpotifyID: func(t *model.Track) string { return t.SpotifyID },
		UpdatedAt: func(t *model.Track) time.Time { return t.UpdatedAt },
	})
	if len(filtered) == 0 {
		return nil
	}

	tracksSpotifyAPI, err := spotifyapi.C.TrackGetAll(ctx, user, filtered)
	if err != nil {
		return nil
	}
	tracksSpotify := utils.SliceMap(tracksSpotifyAPI, func(t spotifyapi.Track) model.Track { return t.ToModel() })

	for i := range tracksSpotify {
		trackDB, ok := utils.SliceFind(tracksDB, func(t *model.Track) bool { return t.Equal(tracksSpotify[i]) })
		if !ok {
			// Track not found
			continue
		}

		tracksSpotify[i].ID = (*trackDB).ID

		// Bring track up to date
		t := tracksSpotify[i]
		if (*trackDB).EqualEntry(t) {
			t = model.Track{ID: t.ID} // Do an empty update to refresh updated_at
		}
		if err := c.track.Update(ctx, t); err != nil {
			return err
		}

		// Bring the track artists up to date
		artistsDB, err := c.artist.GetByTrack(ctx, (*trackDB).ID)
		if err != nil {
			return err
		}

		artistsSpotify := utils.SliceMap(tracksSpotifyAPI[i].Artists, func(a spotifyapi.Artist) model.Artist { return a.ToModel() })

		if err := syncUserData(syncUserDataStruct[model.Artist]{
			DB:     utils.SliceDereference(artistsDB),
			API:    artistsSpotify,
			Equal:  func(a1, a2 model.Artist) bool { return a1.Equal(a2) },
			Get:    func(a model.Artist) (*model.Artist, error) { return c.artist.GetBySpotify(ctx, a.SpotifyID) },
			Create: func(a *model.Artist) error { return c.artist.Create(ctx, a) },
			CreateUserLink: func(a model.Artist) error {
				return c.track.CreateArtist(ctx, &model.TrackArtist{TrackID: (*trackDB).ID, ArtistID: a.ID})
			},
			DeleteUserLink: func(a model.Artist) error {
				return c.track.DeleteArtistByArtistTrack(ctx, model.TrackArtist{TrackID: (*trackDB).ID, ArtistID: a.ID})
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
