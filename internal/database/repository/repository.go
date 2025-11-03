// Package repository interacts with the databank and returns models
package repository

import (
	"context"

	"github.com/topvennie/spotify_organizer/pkg/db"
	"github.com/topvennie/spotify_organizer/pkg/sqlc"
)

type Repository struct {
	db db.DB
}

type contextKey string

const queryKey = contextKey("queries")

func New(db db.DB) *Repository {
	return &Repository{db: db}
}

// Uncomment when creating your first repository
// func (r *Repository) queries(ctx context.Context) *sqlc.Queries {
// 	if q, ok := ctx.Value(queryKey).(*sqlc.Queries); ok {
// 		return q
// 	}
//
// 	return r.db.Queries()
// }

func (r *Repository) WithRollback(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(queryKey).(*sqlc.Queries); ok {
		return fn(ctx)
	}

	return r.db.WithRollback(ctx, func(q *sqlc.Queries) error {
		txCtx := context.WithValue(ctx, queryKey, q)
		return fn(txCtx)
	})
}
