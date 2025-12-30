package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/sqlc"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type History struct {
	repo Repository
}

func (r *Repository) NewHistory() *History {
	return &History{
		repo: *r,
	}
}

func (h *History) GetPreviousPopulated(ctx context.Context, userID int, playedAt time.Time) (*model.History, error) {
	previous, err := h.repo.queries(ctx).HistoryGetPreviousPopulated(ctx, sqlc.HistoryGetPreviousPopulatedParams{
		PlayedAt: toTime(playedAt),
		UserID:   int32(userID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get previous populated %s | %d | %w", playedAt, userID, err)
	}

	previousModel := model.HistoryModel(previous.History)
	previousModel.Track = *model.TrackModel(previous.Track)

	return previousModel, nil
}

func (h *History) GetPopulatedFiltered(ctx context.Context, filter model.HistoryFilter) ([]*model.History, error) {
	skipped := false
	if filter.Skipped != nil {
		skipped = *filter.Skipped
	}

	history, err := h.repo.queries(ctx).HistoryGetPopulatedFiltered(ctx, sqlc.HistoryGetPopulatedFilteredParams{
		Column1:       int32(filter.UserID),
		Limit:         int32(filter.Limit),
		Offset:        int32(filter.Offset),
		Column4:       toTime(filter.Start),
		FilterStart:   !filter.Start.IsZero(),
		Column5:       toTime(filter.End),
		FilterEnd:     !filter.End.IsZero(),
		Column6:       skipped,
		FilterSkipped: filter.Skipped != nil,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get filtered populated history %+v | %w", filter, err)
	}

	return utils.SliceMap(history, func(h sqlc.HistoryGetPopulatedFilteredRow) *model.History {
		history := model.HistoryModel(h.History)
		history.Track = *model.TrackModel(h.Track)

		return history
	}), nil
}

func (h *History) GetSkippedUnknownPopulated(ctx context.Context, userID int) ([]*model.History, error) {
	skippeds, err := h.repo.queries(ctx).HistoryGetSkippedUnknownPopulated(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get skipped unknown populated %d | %w", userID, err)
	}

	return utils.SliceMap(skippeds, func(s sqlc.HistoryGetSkippedUnknownPopulatedRow) *model.History {
		h := model.HistoryModel(s.History)
		h.Track = *model.TrackModel(s.Track)

		return h
	}), nil
}

func (h *History) Create(ctx context.Context, history *model.History) error {
	id, err := h.repo.queries(ctx).HistoryCreate(ctx, sqlc.HistoryCreateParams{
		UserID:     int32(history.UserID),
		TrackID:    int32(history.TrackID),
		PlayedAt:   toTime(history.PlayedAt),
		Skipped:    toBool(history.Skipped),
		AlbumID:    toInt(history.AlbumID),
		ArtistID:   toInt(history.ArtistID),
		PlaylistID: toInt(history.PlaylistID),
		ShowID:     toInt(history.ShowID),
	})
	if err != nil {
		return fmt.Errorf("create history %+v | %w", *history, err)
	}

	history.ID = int(id)

	return nil
}

// CreateBatch let's you create a batch of history entries
// It will panic if the skipped field is nil!!!
func (h *History) CreateBatch(ctx context.Context, histories []model.History) error {
	userIDs := make([]int32, 0, len(histories))
	trackIDs := make([]int32, 0, len(histories))
	playedAts := make([]pgtype.Timestamptz, 0, len(histories))
	skippeds := make([]bool, 0, len(histories))

	for i := range histories {
		userIDs = append(userIDs, int32(histories[i].UserID))
		trackIDs = append(trackIDs, int32(histories[i].TrackID))
		playedAts = append(playedAts, toTime(histories[i].PlayedAt))
		skippeds = append(skippeds, *histories[i].Skipped)
	}

	if err := h.repo.queries(ctx).HistoryCreateBatch(ctx, sqlc.HistoryCreateBatchParams{
		Column1: userIDs,
		Column2: trackIDs,
		Column3: playedAts,
		Column4: skippeds,
	}); err != nil {
		return fmt.Errorf("create history batch %w", err)
	}

	return nil
}

func (h *History) Update(ctx context.Context, history model.History) error {
	if err := h.repo.queries(ctx).HistoryUpdate(ctx, sqlc.HistoryUpdateParams{
		ID:       int32(history.ID),
		PlayedAt: toTime(history.PlayedAt),
		Skipped:  toBool(history.Skipped),
	}); err != nil {
		return fmt.Errorf("update history %+v | %w", history, err)
	}

	return nil
}

func (h *History) DeleteOlder(ctx context.Context, userID int, playedAt time.Time) error {
	if err := h.repo.queries(ctx).HistoryDeleteUserOlder(ctx, sqlc.HistoryDeleteUserOlderParams{
		UserID:   int32(userID),
		PlayedAt: toTime(playedAt),
	}); err != nil {
		return fmt.Errorf("delete history for user %d by time %s", userID, playedAt)
	}

	return nil
}
