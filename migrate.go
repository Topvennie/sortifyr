package main

import (
	"embed"

	"github.com/topvennie/spotify_organizer/pkg/config"
	"github.com/topvennie/spotify_organizer/pkg/db"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	// setup database
	db, err := db.NewPSQL()
	if err != nil {
		panic(err)
	}
	conn := stdlib.OpenDBFromPool(db.Pool())

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(conn, "db/migrations"); err != nil {
		panic(err)
	}

	// run app
}
