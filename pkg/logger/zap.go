// Package logger initiates a zap logger
package logger

import (
	"github.com/topvennie/spotify_organizer/pkg/config"
	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	var logger *zap.Logger

	if config.IsDev() {
		logger = zap.Must(zap.NewDevelopment(zap.AddStacktrace(zap.WarnLevel)))
	} else {
		cfg := zap.NewProductionConfig()
		cfg.Level.SetLevel(zap.WarnLevel)
		logger = zap.Must(zap.NewProduction())
	}

	env := config.GetDefaultString("app.env", "development")
	logger = logger.With(zap.String("env", env))

	return logger, nil
}
