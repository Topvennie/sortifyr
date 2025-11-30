// Package storage connects with a file / image storage
package storage

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/minio"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/topvennie/spotify_organizer/pkg/config"
)

var S fiber.Storage

func New(pool *pgxpool.Pool) error {
	provider := config.GetDefaultString("storage.provider", "minio")

	switch provider {

	case "minio":
		S = minio.New(minio.Config{
			Bucket:   config.GetDefaultString("minio.bucket", "spotify"),
			Endpoint: config.GetDefaultString("minio.endpoint", "minio:9000"),
			Secure:   config.GetDefaultBool("minio.secure", false),
			Credentials: minio.Credentials{
				AccessKeyID:     config.GetDefaultString("minio.username", "minio"),
				SecretAccessKey: config.GetDefaultString("minio.password", "miniominio"),
			},
		})

	case "postgres":
		S = postgres.New(postgres.Config{
			DB:         pool,
			Table:      "spotify_organizer_files",
			Reset:      false,
			GCInterval: 10 * time.Second,
		})

	default:
		return fmt.Errorf("unsupported storage provider %s", provider)
	}

	return nil
}
