package spotifyapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/topvennie/sortifyr/internal/database/model"
)

func (c *client) AlbumGet(ctx context.Context, user model.User, spotifyID string) (Album, error) {
	var resp Album

	if err := c.request(ctx, user, http.MethodGet, "albums/"+spotifyID, http.NoBody, &resp); err != nil {
		return Album{}, fmt.Errorf("get album %s | %w", spotifyID, err)
	}

	return resp, nil
}

type albumUserResponse struct {
	Total int     `json:"total"`
	Items []Album `json:"items"`
}

func (c *client) AlbumGetUser(ctx context.Context, user model.User) ([]Album, error) {
	albums := make([]Album, 0)

	total := 51
	limit := 50

	for i := 0; i < total; i += limit {
		var resp albumUserResponse

		if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("me/albums?offset=%d&limit=%d", i, limit), http.NoBody, &resp); err != nil {
			return nil, fmt.Errorf("get albums with limit %d and offset %d | %w", limit, i, err)
		}

		albums = append(albums, resp.Items...)
		total = resp.Total
	}

	return albums, nil
}

type albumAllResponse struct {
	Albums []Album `json:"albums"`
}

func (c *client) AlbumGetAll(ctx context.Context, user model.User, albumIDs []string) ([]Album, error) {
	albums := make([]Album, 0, len(albumIDs))

	limit := 20

	for i := 0; i < len(albumIDs); i += limit {
		var resp albumAllResponse

		url := "albums?ids=" + strings.Join(albumIDs[i:min(len(albumIDs), i+limit)], ",")
		if err := c.request(ctx, user, http.MethodGet, url, http.NoBody, &resp); err != nil {
			return nil, err
		}

		albums = append(albums, resp.Albums...)
	}

	return albums, nil
}

type albumTrackResponse struct {
	Total int     `json:"total"`
	Items []Track `json:"items"`
}

func (c *client) AlbumGetTrackAll(ctx context.Context, user model.User, spotifyID string) ([]Track, error) {
	tracks := make([]Track, 0)

	total := 51
	limit := 50

	for i := 0; i < total; i += limit {
		var resp albumTrackResponse

		if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("albums/%s/tracks?offset=%d&limit=%d", spotifyID, i, limit), http.NoBody, &resp); err != nil {
			return nil, err
		}

		tracks = append(tracks, resp.Items...)
		total = resp.Total
	}

	return tracks, nil
}
