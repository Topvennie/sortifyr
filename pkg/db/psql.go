package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/topvennie/spotify_organizer/pkg/config"
	"github.com/topvennie/spotify_organizer/pkg/sqlc"
)

type psql struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// Interface compliance
var _ DB = (*psql)(nil)

func NewPSQL() (DB, error) {
	pgConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	pgConfig.ConnConfig.Host = config.GetDefaultString("db.host", "localhost")
	pgConfig.ConnConfig.Port = config.GetDefaultUint16("db.port", 5432)
	pgConfig.ConnConfig.Database = config.GetDefaultString("db.database", "spotify_organizer")
	pgConfig.ConnConfig.User = config.GetDefaultString("db.user", "postgres")
	pgConfig.ConnConfig.Password = config.GetDefaultString("db.password", "postgres")

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.TODO()); err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	return &psql{pool: pool, queries: queries}, nil
}

func (p *psql) WithRollback(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		errRollback := tx.Rollback(ctx)
		if !errors.Is(err, pgx.ErrTxClosed) {
			err = errRollback
		}
	}()

	queries := sqlc.New(tx)

	if err := fn(queries); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *psql) Pool() *pgxpool.Pool {
	return p.pool
}

func (p *psql) Queries() *sqlc.Queries {
	return p.queries
}
