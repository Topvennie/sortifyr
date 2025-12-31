package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
	"github.com/topvennie/sortifyr/internal/server/dto"
	"github.com/topvennie/sortifyr/pkg/utils"
	"go.uber.org/zap"
)

type Track struct {
	service Service

	history repository.History
	track   repository.Track
}

func (s *Service) NewTrack() *Track {
	return &Track{
		service: *s,
		history: *s.repo.NewHistory(),
		track:   *s.repo.NewTrack(),
	}
}

func (t *Track) GetHistory(ctx context.Context, filter dto.HistoryFilter) ([]dto.History, error) {
	filterModel := filter.ToModel()
	if filterModel.Skipped == nil {
		tmp := false
		filterModel.PlayCountSkipped = &tmp
	}

	history, err := t.history.GetPopulatedFiltered(ctx, *filterModel)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(history, func(h *model.History) dto.History { return dto.HistoryDTO(&h.Track, h) }), nil
}

func (t *Track) GetAdded(ctx context.Context, filter dto.TrackFilter) ([]dto.TrackAdded, error) {
	tracks, err := t.track.GetCreatedFiltered(ctx, *filter.ToModel())
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(tracks, dto.TrackAddedDTO), nil
}

func (t *Track) GetDeleted(ctx context.Context, filter dto.TrackFilter) ([]dto.TrackDeleted, error) {
	tracks, err := t.track.GetDeletedFiltered(ctx, *filter.ToModel())
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(tracks, dto.TrackDeletedDTO), nil
}
