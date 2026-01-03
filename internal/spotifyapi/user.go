package spotifyapi

import (
	"context"
	"net/http"

	"github.com/topvennie/sortifyr/internal/database/model"
)

type userResponse struct {
	DisplayName string `json:"display_name"`
}

func (c *client) UserGet(ctx context.Context, user, spotifyUser model.User) (model.User, error) {
	var resp userResponse

	if err := c.request(ctx, user, http.MethodGet, "users/"+spotifyUser.UID, http.NoBody, &resp); err != nil {
		return model.User{}, err
	}

	spotifyUser.DisplayName = resp.DisplayName

	return spotifyUser, nil
}
