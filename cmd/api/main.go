package main

import (
	"fmt"

	"github.com/topvennie/sortifyr/internal/database/repository"
	"github.com/topvennie/sortifyr/internal/generator"
	"github.com/topvennie/sortifyr/internal/server"
	"github.com/topvennie/sortifyr/internal/server/service"
	"github.com/topvennie/sortifyr/internal/spotify"
	"github.com/topvennie/sortifyr/internal/task"
	"github.com/topvennie/sortifyr/pkg/config"
	"github.com/topvennie/sortifyr/pkg/db"
	"github.com/topvennie/sortifyr/pkg/logger"
	"github.com/topvennie/sortifyr/pkg/redis"
	"github.com/topvennie/sortifyr/pkg/storage"
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

	if err := task.Init(*repo); err != nil {
		zap.S().Fatalf("Failed to init the task package %v", err)
	}

	if err := spotify.Init(*repo); err != nil {
		zap.S().Fatalf("Failed to init the spotify package %v", err)
	}

	generator.Init(*repo)

	api := server.New(*service, db.Pool())

	zap.S().Infof("Server is running on %s", api.Addr)
	if err := api.Listen(api.Addr); err != nil {
		zap.S().Fatalf("Failure while running the server %v", err)
	}
}
