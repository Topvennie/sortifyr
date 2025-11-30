// Package redis connects to the redis db
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/topvennie/spotify_organizer/pkg/config"
)

var (
	C      *redis.Client
	ErrNil = redis.Nil
)

func New() error {
	URL := config.GetDefaultString("redis.url", "redis://default@redis:6379")

	options, err := redis.ParseURL(URL)
	if err != nil {
		return err
	}

	C = redis.NewClient(options)
	ctx := context.Background()
	_, err = C.Ping(ctx).Result()
	return err
}
