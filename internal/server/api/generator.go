package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/sortifyr/internal/server/dto"
	"github.com/topvennie/sortifyr/internal/server/service"
)

type Generator struct {
	router fiber.Router

	generator service.Generator
}

func NewGenerator(router fiber.Router, service service.Service) *Generator {
	api := &Generator{
		router:    router.Group("/generator"),
		generator: *service.NewGenerator(),
	}

	api.createRoutes()

	return api
}

func (g *Generator) createRoutes() {
	g.router.Post("/generate", g.generate)
}

func (g *Generator) generate(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var generator dto.Generator
	if err := c.BodyParser(&generator); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(generator); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tracks, err := g.generator.Generate(c.Context(), userID, generator)
	if err != nil {
		return err
	}

	return c.JSON(tracks)
}
