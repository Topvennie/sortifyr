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
	g.router.Get("/", g.getAll)
	g.router.Post("/preview", g.preview)
	g.router.Put("/", g.create)
	g.router.Post("/:id", g.edit)
}

func (g *Generator) getAll(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	generators, err := g.generator.GetByUser(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(generators)
}

func (g *Generator) preview(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var generatorParams dto.GeneratorParams
	if err := c.BodyParser(&generatorParams); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(generatorParams); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tracks, err := g.generator.Preview(c.Context(), userID, generatorParams)
	if err != nil {
		return err
	}

	return c.JSON(tracks)
}

func (g *Generator) create(c *fiber.Ctx) error {
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
	generator.ID = 0

	generator, err := g.generator.Create(c.Context(), userID, generator)
	if err != nil {
		return err
	}

	return c.JSON(generator)
}

func (g *Generator) edit(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	id, err := c.ParamsInt(":id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var generator dto.Generator
	if err := c.BodyParser(&generator); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(generator); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	generator.ID = id

	generator, err = g.generator.Edit(c.Context(), userID, generator)
	if err != nil {
		return err
	}

	return c.JSON(generator)
}
