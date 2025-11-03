// Package server starts the server
package server

import (
	"fmt"

	"github.com/topvennie/spotify_organizer/internal/server/service"
	"github.com/topvennie/spotify_organizer/pkg/config"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	redisstore "github.com/gofiber/storage/redis/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

type Server struct {
	*fiber.App
	Addr string
}

func New(_ service.Service, pool *pgxpool.Pool) *Server {
	// Construct app
	app := fiber.New(fiber.Config{
		BodyLimit:         20 * 1024 * 1024,
		ReadBufferSize:    8096,
		StreamRequestBody: true,
	})

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
	}))
	if config.IsDev() {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:3000",
			AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
			AllowCredentials: true,
		}))
	}

	// Session storage
	var sessionStore fiber.Storage
	if !config.IsDev() {
		sessionStore = redisstore.New(redisstore.Config{
			URL: config.GetDefaultString("redis.url", ""),
		})
	} else {
		sessionStore = postgres.New(postgres.Config{
			DB: pool,
		})
	}

	goth_fiber.SessionStore = session.New(session.Config{
		KeyLookup:      fmt.Sprintf("cookie:%s_session_id", config.GetString("app.name")),
		CookieHTTPOnly: true,
		Storage:        sessionStore,
		CookieSecure:   !config.IsDev(),
	})

	// Register routes
	_ = app.Group("/api")

	// Add new routes here

	// Static files if served in production
	if !config.IsDev() {
		app.Static("/", "./public")
		// Fallback for SPA to handle
		app.Static("*", "./public/index.html")
	}

	// Fallback
	app.All("/api*", func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	port := config.GetDefaultInt("server.port", 8000)
	host := config.GetDefaultString("server.host", "0.0.0.0")

	srv := &Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		App:  app,
	}

	return srv
}
