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

func (c *client) TrackGet(ctx context.Context, user model.User, spotifyID string) (Track, error) {
	var resp Track

	if err := c.request(ctx, user, http.MethodGet, "tracks/"+spotifyID, http.NoBody, &resp); err != nil {
		return Track{}, fmt.Errorf("get track %s | %w", spotifyID, err)
	}

	return resp, nil
}

type trackAllResponse struct {
	Tracks []Track `json:"tracks"`
}

func (c *client) TrackGetAll(ctx context.Context, user model.User, trackIDs []string) ([]Track, error) {
	wg := concurrent.NewLimitedWaitGroup(12)

	var mu sync.Mutex
	var errs []error

	tracks := make([]Track, 0, len(trackIDs))
	limit := 50

	for i := 0; i < len(trackIDs); i += limit {
		wg.Go(func() {
			var resp trackAllResponse

			url := "tracks?ids=" + strings.Join(trackIDs[i:min(len(trackIDs), i+limit)], ",")
			if err := c.request(ctx, user, http.MethodGet, url, http.NoBody, &resp); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			tracks = append(tracks, resp.Tracks...)
			mu.Unlock()
		})
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return tracks, nil
}
