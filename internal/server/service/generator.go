package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
	"github.com/topvennie/sortifyr/internal/generator"
	"github.com/topvennie/sortifyr/internal/server/dto"
	"github.com/topvennie/sortifyr/pkg/utils"
	"go.uber.org/zap"
)

type Generator struct {
	generator repository.Generator
	user      repository.User
}

func (s *Service) NewGenerator() *Generator {
	return &Generator{
		generator: *s.repo.NewGenerator(),
		user:      *s.repo.NewUser(),
	}
}

func (g *Generator) GetByUser(ctx context.Context, userID int) ([]dto.Generator, error) {
	gens, err := g.generator.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(gens, dto.GeneratorDTO), nil
}

func (g *Generator) Preview(ctx context.Context, userID int, params dto.GeneratorParams) ([]dto.Track, error) {
	gen := model.Generator{
		UserID: userID,
		Params: params.ToModel(),
	}
	tracks, err := generator.G.Generate(ctx, &gen)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(tracks, func(t model.Track) dto.Track { return dto.TrackDTO(&t) }), nil
}

func (g *Generator) Refresh(ctx context.Context, userID, genID int) error {
	user, err := g.user.GetByID(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return fiber.ErrUnauthorized
	}

	gen, err := g.generator.Get(ctx, genID)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	if gen == nil {
		return fiber.ErrNotFound
	}

	if gen.PlaylistID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "generator does not have a playlist")
	}

	if gen.UserID != userID {
		return fiber.ErrForbidden
	}

	if err := generator.G.Refresh(*user, *gen); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (g *Generator) Create(ctx context.Context, userID int, genSave dto.GeneratorSave) (dto.Generator, error) {
	gen := genSave.ToModel(userID)

	if !genSave.CreatePlaylist {
		gen.Maintained = false
	}
	if !genSave.Maintained {
		gen.Interval = 0
	}
	if gen.Maintained && gen.Interval < 24*time.Hour {
		return dto.Generator{}, fiber.NewError(fiber.StatusBadRequest, "interval needs to be > 0 if maintained")
	}

	zap.S().Debug(genSave)
	zap.S().Debug(*gen)

	if err := generator.G.Create(ctx, gen, genSave.CreatePlaylist); err != nil {
		zap.S().Error(err)
		return dto.Generator{}, fiber.ErrInternalServerError
	}

	zap.S().Debug(*gen)

	return dto.GeneratorDTO(gen), nil
}

func (g *Generator) Update(ctx context.Context, userID int, genSave dto.GeneratorSave) (dto.Generator, error) {
	oldGen, err := g.generator.Get(ctx, genSave.ID)
	if err != nil {
		zap.S().Error(err)
		return dto.Generator{}, fiber.ErrInternalServerError
	}
	if oldGen == nil {
		return dto.Generator{}, fiber.ErrNotFound
	}
	if oldGen.UserID != userID {
		return dto.Generator{}, fiber.ErrForbidden
	}

	gen := genSave.ToModel(userID)

	if !genSave.CreatePlaylist {
		gen.Maintained = false
	}
	if !genSave.Maintained {
		gen.Interval = 0
	}

	if err := generator.G.Update(ctx, gen, genSave.CreatePlaylist); err != nil {
		zap.S().Error(err)
		return dto.Generator{}, fiber.ErrInternalServerError
	}

	return dto.GeneratorDTO(gen), nil
}
