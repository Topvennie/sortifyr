// Package server starts the server
package server

import (
	"fmt"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shareed2k/goth_fiber"
	routers "github.com/topvennie/sortifyr/internal/server/api"
	"github.com/topvennie/sortifyr/internal/server/middlewares"
	"github.com/topvennie/sortifyr/internal/server/service"
	"github.com/topvennie/sortifyr/pkg/config"

	"go.uber.org/zap"
)

type Server struct {
	*fiber.App
	Addr string
}

func New(service service.Service, pool *pgxpool.Pool) *Server {
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
	sessionStore := postgres.New(postgres.Config{
		DB: pool,
	})

	goth_fiber.SessionStore = session.New(session.Config{
		KeyLookup:      fmt.Sprintf("cookie:%s_session_id", config.GetDefaultString("app.name", "sortifyr")),
		CookieHTTPOnly: true,
		Storage:        sessionStore,
		CookieSecure:   !config.IsDev(),
	})

	// Register routes
	api := app.Group("/api")

	routers.NewAuth(api, service)

	protectedAPI := api.Use(middlewares.ProtectedRoute)

	routers.NewUser(protectedAPI, service)
	routers.NewPlaylist(protectedAPI, service)
	routers.NewDirectory(protectedAPI, service)
	routers.NewLink(protectedAPI, service)
	routers.NewTask(protectedAPI, service)
	routers.NewTrack(protectedAPI, service)
	routers.NewGenerator(protectedAPI, service)
	routers.NewSetting(protectedAPI, service)

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
