// Package redis connects to the redis db
package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/topvennie/spotify_organizer/pkg/config"
)

var (
	C      *redis.Client
	ErrNil = redis.Nil
)

func New() error {
	URL := config.GetDefaultString("redis.url", "")
	if URL == "" {
		return errors.New("no redis url configured")
	}

	options, err := redis.ParseURL(URL)
	if err != nil {
		return err
	}

	C = redis.NewClient(options)
	ctx := context.Background()
	_, err = C.Ping(ctx).Result()
	return err
}
