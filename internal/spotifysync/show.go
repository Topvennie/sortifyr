package spotifysync

import (
	"context"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/pkg/utils"
)

// showSync will syncronize the user's saved shows
func (c *client) showSync(ctx context.Context, user model.User) error {
	showsDB, err := c.show.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	showsSpotifyAPI, err := spotifyapi.C.ShowGetUser(ctx, user)
	if err != nil {
		return err
	}
	showsSpotify := utils.SliceMap(showsSpotifyAPI, func(s spotifyapi.Show) model.Show { return s.ToModel() })

	return syncUserData(syncUserDataStruct[model.Show]{
		DB:     utils.SliceDereference(showsDB),
		API:    showsSpotify,
		Equal:  func(s1, s2 model.Show) bool { return s1.Equal(s2) },
		Get:    func(s model.Show) (*model.Show, error) { return c.show.GetBySpotify(ctx, s.SpotifyID) },
		Create: func(s *model.Show) error { return c.show.Create(ctx, s) },
		CreateUserLink: func(s model.Show) error {
			return c.show.CreateUser(ctx, &model.ShowUser{ShowID: s.ID, UserID: user.ID})
		},
		DeleteUserLink: func(s model.Show) error {
			return c.show.DeleteUserByUserShow(ctx, model.ShowUser{ShowID: s.ID, UserID: user.ID})
		},
	})
}

// showUpdate updates local show instances to match the spotify data
// It updates all shows, regardless of the user given.
// However the given user's access token is used.
func (c *client) showUpdate(ctx context.Context, user model.User) error {
	showsDB, err := c.show.GetAll(ctx)
	if err != nil {
		return err
	}

	filtered := filterSpotify(filterSpotifyStruct[*model.Show]{
		Items:     showsDB,
		Frequency: 24,
		SpotifyID: func(s *model.Show) string { return s.SpotifyID },
		UpdatedAt: func(s *model.Show) time.Time { return s.UpdatedAt },
	})
	if len(filtered) == 0 {
		return nil
	}

	showsSpotifyAPI, err := spotifyapi.C.ShowGetAll(ctx, user, utils.SliceMap(showsDB, func(s *model.Show) string { return s.SpotifyID }))
	if err != nil {
		return err
	}
	showsSpotify := utils.SliceMap(showsSpotifyAPI, func(s spotifyapi.Show) model.Show { return s.ToModel() })

	for i := range showsSpotify {
		showDB, ok := utils.SliceFind(showsDB, func(s *model.Show) bool { return s.Equal(showsSpotify[i]) })
		if !ok {
			// Show not found
			continue
		}

		showsSpotify[i].ID = (*showDB).ID

		// Bring the show up to date
		s := showsSpotify[i]
		if (*showDB).EqualEntry(s) {
			s = model.Show{ID: s.ID}
		}
		if err := c.show.Update(ctx, showsSpotify[i]); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) showCoverSync(ctx context.Context, user model.User) error {
	shows, err := c.show.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	return c.syncCover(ctx, utils.SliceMap(shows, func(s *model.Show) syncCoverStruct {
		return syncCoverStruct{
			CoverURL: s.CoverURL,
			CoverID:  s.CoverID,
			Update: func(newID string) error {
				s.CoverID = newID
				return c.show.Update(ctx, *s)
			},
		}
	}))
}
