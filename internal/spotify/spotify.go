package spotify

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/internal/database/repository"
	"github.com/topvennie/spotify_organizer/internal/task"
	"github.com/topvennie/spotify_organizer/pkg/config"
	"github.com/topvennie/spotify_organizer/pkg/redis"
	"go.uber.org/zap"
)

var ErrUnauthorized = errors.New("access and refresh token expired")

type client struct {
	directory repository.Directory
	link      repository.Link
	playlist  repository.Playlist
	track     repository.Track
	user      repository.User

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
		directory:    *repo.NewDirectory(),
		link:         *repo.NewLink(),
		playlist:     *repo.NewPlaylist(),
		track:        *repo.NewTrack(),
		user:         *repo.NewUser(),
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	if err := C.taskRegister(); err != nil {
		return err
	}

	return nil
}

func (c *client) NewUser(ctx context.Context, user model.User, accessToken, refreshToken string, expiresIn time.Duration) error {
	if _, err := redis.C.Set(ctx, accessKey(user), accessToken, expiresIn).Result(); err != nil {
		return fmt.Errorf("set access token %w", err)
	}

	if _, err := redis.C.Set(ctx, refreshKey(user), refreshToken, 0).Result(); err != nil {
		return fmt.Errorf("set refresh token %w", err)
	}

	if err := task.Manager.RunByUID(taskPlaylistUID, user); err != nil {
		zap.S().Fatalf("%v", err)
	}

	if err := task.Manager.RunByUID(taskTrackUID, user); err != nil {
		zap.S().Fatalf("%v", err)
	}

	if err := task.Manager.RunByUID(taskUserUID, user); err != nil {
		zap.S().Fatalf("%v", err)
	}

	return nil
}
