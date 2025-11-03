// Package db connects with the databank
package db

import (
	"context"

	"github.com/topvennie/spotify_organizer/pkg/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	WithRollback(ctx context.Context, fn func(q *sqlc.Queries) error) error
	Pool() *pgxpool.Pool
	Queries() *sqlc.Queries
}
