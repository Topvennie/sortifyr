package main

import (
	"fmt"

	"github.com/topvennie/spotify_organizer/internal/database/repository"
	"github.com/topvennie/spotify_organizer/internal/server"
	"github.com/topvennie/spotify_organizer/internal/server/service"
	"github.com/topvennie/spotify_organizer/internal/spotify"
	"github.com/topvennie/spotify_organizer/pkg/config"
	"github.com/topvennie/spotify_organizer/pkg/db"
	"github.com/topvennie/spotify_organizer/pkg/logger"
	"github.com/topvennie/spotify_organizer/pkg/redis"
	"github.com/topvennie/spotify_organizer/pkg/storage"
	"go.uber.org/zap"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	zapLogger, err := logger.New()
	if err != nil {
		panic(fmt.Errorf("zap logger initialization failed: %w", err))
	}
	zap.ReplaceGlobals(zapLogger)

	db, err := db.NewPSQL()
	if err != nil {
		zap.S().Fatalf("Unable to connect to database %v", err)
	}

	if err = storage.New(db.Pool()); err != nil {
		zap.S().Fatalf("Failed to create storage %v", err)
	}

	if err = redis.New(); err != nil {
		zap.S().Fatalf("Failed to connect to redis %v", err)
	}

	repo := repository.New(db)
	service := service.New(*repo)

	if err := spotify.Init(*repo); err != nil {
		zap.S().Fatalf("Failed to init the spotify package %v", err)
	}

	api := server.New(*service, db.Pool())

	zap.S().Infof("Server is running on %s", api.Addr)
	if err := api.Listen(api.Addr); err != nil {
		zap.S().Fatalf("Failure while running the server %v", err)
	}
}
