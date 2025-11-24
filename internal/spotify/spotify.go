package spotify

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/internal/database/repository"
	"github.com/topvennie/spotify_organizer/pkg/config"
	"github.com/topvennie/spotify_organizer/pkg/redis"
)

type client struct {
	playlist repository.Playlist
	user     repository.User

	clientID     string
	clientSecret string
}

var C *client

func Init(repo repository.Repository) error {
	clientID := config.GetString("auth.spotify.client.id")
	clientSecret := config.GetString("auth.spotify.client.secret")

	if clientID == "" || clientSecret == "" {
		return errors.New("client id or client secret not set")
	}

	C = &client{
		playlist:     *repo.NewPlaylist(),
		user:         *repo.NewUser(),
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	return nil
}

func (c *client) NewUser(ctx context.Context, user model.User, accessToken, refreshToken string, expiresIn time.Duration) error {
	if _, err := redis.C.Set(ctx, accessKey(user.UID), accessToken, expiresIn).Result(); err != nil {
		return fmt.Errorf("set access token %w", err)
	}

	if _, err := redis.C.Set(ctx, refreshKey(user.UID), refreshToken, 0).Result(); err != nil {
		return fmt.Errorf("set refresh token %w", err)
	}

	userFull, err := c.getUser(ctx, user.UID)
	if err != nil {
		return err
	}

	if !user.Equal(userFull) {
		if err := c.user.Update(ctx, userFull); err != nil {
			return err
		}
	}

	return nil
}

func accessKey(uid string) string {
	return uid + ":spotify:access_token"
}

func refreshKey(uid string) string {
	return uid + ":spotify:refresh_token"
}
