package spotify

import (
	"context"
	"fmt"

	"github.com/topvennie/spotify_organizer/internal/database/model"
)

type userResponse struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *client) getUser(ctx context.Context, uid string) (model.User, error) {
	user, err := c.user.GetByUID(ctx, uid)
	if err != nil {
		return model.User{}, err
	}
	if user == nil {
		return model.User{}, fmt.Errorf("no user found for uid %s", uid)
	}

	var resp userResponse

	if err := c.request(ctx, uid, "me", &resp); err != nil {
		return model.User{}, nil
	}

	user.DisplayName = resp.DisplayName

	return *user, nil
}
