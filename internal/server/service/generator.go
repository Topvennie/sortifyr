package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/generator"
	"github.com/topvennie/sortifyr/internal/server/dto"
	"github.com/topvennie/sortifyr/pkg/utils"
	"go.uber.org/zap"
)

type Generator struct{}

func (s *Service) NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Preview(ctx context.Context, userID int, params dto.GeneratorParams) ([]dto.Track, error) {
	gen := model.Generator{
		UserID: userID,
		Params: params.ToModel(),
	}
	tracks, err := generator.G.Generate(ctx, gen)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(tracks, func(t model.Track) dto.Track { return dto.TrackDTO(&t) }), nil
}
