// Package redis connects to the redis db
package redis

import (
	"context"
	"errors"

	"github.com/topvennie/spotify_organizer/pkg/config"
	"github.com/redis/go-redis/v9"
)

var (
	C      *redis.Client
	ErrNil = redis.Nil
)

func New() error {
	url := config.GetDefaultString("redis.url", "")
	if url == "" {
		return errors.New("no redis url configured")
	}

	options, err := redis.ParseURL(url)
	if err != nil {
		return err
	}

	C = redis.NewClient(options)
	ctx := context.Background()
	_, err = C.Ping(ctx).Result()
	return err
}
