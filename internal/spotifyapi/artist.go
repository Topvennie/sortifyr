package spotifyapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/concurrent"
)

func (c *client) ArtistGet(ctx context.Context, user model.User, spotifyID string) (Artist, error) {
	var resp Artist

	if err := c.request(ctx, user, http.MethodGet, "artists/"+spotifyID, http.NoBody, &resp); err != nil {
		return Artist{}, fmt.Errorf("get artist %s | %w", spotifyID, err)
	}

	return resp, nil
}

type artistAllResponse struct {
	Artists []Artist `json:"artists"`
}

func (c *client) ArtistGetAll(ctx context.Context, user model.User, artistIDs []string) ([]Artist, error) {
	wg := concurrent.NewLimitedWaitGroup(12)

	var mu sync.Mutex
	var errs []error

	artists := make([]Artist, 0, len(artistIDs))
	limit := 50

	for i := 0; i < len(artistIDs); i += limit {
		wg.Go(func() {
			var resp artistAllResponse

			url := "artists?ids=" + strings.Join(artistIDs[i:min(len(artistIDs), i+limit)], ",")
			if err := c.request(ctx, user, http.MethodGet, url, http.NoBody, &resp); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			artists = append(artists, resp.Artists...)
			mu.Unlock()
		})
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return artists, nil
}
