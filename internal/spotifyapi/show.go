package spotifyapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/topvennie/sortifyr/internal/database/model"
)

func (c *client) ShowGet(ctx context.Context, user model.User, spotifyID string) (Show, error) {
	var resp Show

	if err := c.request(ctx, user, http.MethodGet, "shows/"+spotifyID, http.NoBody, &resp); err != nil {
		return Show{}, fmt.Errorf("get show %s | %w", spotifyID, err)
	}

	return resp, nil
}

type showUserResponse struct {
	Total int    `json:"total"`
	Items []Show `json:"items"`
}

func (c *client) ShowGetUser(ctx context.Context, user model.User) ([]Show, error) {
	shows := make([]Show, 0)

	total := 51
	limit := 50

	for i := 0; i < total; i += limit {
		var resp showUserResponse

		if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("me/shows?offset=%d&limit=%d", i, limit), http.NoBody, &resp); err != nil {
			return nil, fmt.Errorf("get shows with limit %d and offset %d | %w", limit, i, err)
		}

		shows = append(shows, resp.Items...)
		total = resp.Total
	}

	return shows, nil
}

type showAllResponse struct {
	Shows []Show `json:"shows"`
}

func (c *client) ShowGetAll(ctx context.Context, user model.User, showsIDs []string) ([]Show, error) {
	shows := make([]Show, 0, len(showsIDs))

	limit := 50

	for i := 0; i < len(showsIDs); i += limit {
		var resp showAllResponse

		url := "shows?ids=" + strings.Join(showsIDs[i:min(len(showsIDs), i+limit)], ",")
		if err := c.request(ctx, user, http.MethodGet, url, http.NoBody, &resp); err != nil {
			return nil, err
		}

		shows = append(shows, resp.Shows...)
	}

	return shows, nil
}
