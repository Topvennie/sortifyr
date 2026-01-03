package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/sqlc"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type Track struct {
	repo Repository
}

func (r *Repository) NewTrack() *Track {
	return &Track{
		repo: *r,
	}
}

func (t *Track) GetAll(ctx context.Context) ([]*model.Track, error) {
	tracks, err := t.repo.queries(ctx).TrackGetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all track %w", err)
	}

	return utils.SliceMap(tracks, model.TrackModel), nil
}

func (t *Track) GetAllByID(ctx context.Context, trackIDs []int) ([]*model.Track, error) {
	tracks, err := t.repo.queries(ctx).TrackGetAllById(ctx, utils.SliceMap(trackIDs, func(t int) int32 { return int32(t) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all tracks by id %+v | %w", trackIDs, err)
	}

	return utils.SliceMap(tracks, model.TrackModel), nil
}

func (t *Track) GetBySpotify(ctx context.Context, spotifyID string) (*model.Track, error) {
	track, err := t.repo.queries(ctx).TrackGetBySpotify(ctx, spotifyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get track by spotify id %s | %w", spotifyID, err)
	}

	return model.TrackModel(track), nil
}

func (t *Track) GetAllBySpotify(ctx context.Context, spotifyIDs []string) ([]*model.Track, error) {
	tracks, err := t.repo.queries(ctx).TrackGetAllBySpotify(ctx, spotifyIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all tracks by spotify id %w", err)
	}

	return utils.SliceMap(tracks, model.TrackModel), nil
}

func (t *Track) GetByName(ctx context.Context, name string) ([]*model.Track, error) {
	tracks, err := t.repo.queries(ctx).TrackGetByName(ctx, toString(name))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get tracks by name %s | %w", name, err)
	}

	return utils.SliceMap(tracks, model.TrackModel), nil
}

func (t *Track) GetByPlaylist(ctx context.Context, playlistID int) ([]*model.Track, error) {
	tracks, err := t.repo.queries(ctx).TrackGetByPlaylist(ctx, int32(playlistID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get tracks by playlist %d | %w", playlistID, err)
	}

	return utils.SliceMap(tracks, model.TrackModel), nil
}

func (t *Track) GetCreatedFiltered(ctx context.Context, filter model.TrackFilter) ([]*model.Track, error) {
	params := sqlc.TrackGetCreatedFilteredPopulatedParams{
		Column3:          int32(filter.UserID),
		Column4:          int32(filter.PlaylistID),
		FilterPlaylistID: filter.PlaylistID != 0,
		Limit:            int32(filter.Limit),
		Offset:           int32(filter.Offset),
	}

	tracks, err := t.repo.queries(ctx).TrackGetCreatedFilteredPopulated(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get filtered created tracks %+v | %w", filter, err)
	}

	return utils.SliceMap(tracks, func(t sqlc.TrackGetCreatedFilteredPopulatedRow) *model.Track {
		track := model.TrackModel(t.Track)
		track.Playlist = *model.PlaylistModel(t.Playlist)
		track.Playlist.Owner = *model.UserModel(t.User)
		track.CreatedAt = t.PlaylistTrack.CreatedAt.Time

		return track
	}), nil
}

func (t *Track) GetDeletedFiltered(ctx context.Context, filter model.TrackFilter) ([]*model.Track, error) {
	params := sqlc.TrackGetDeletedFilteredPopulatedParams{
		Column3:          int32(filter.UserID),
		Column4:          int32(filter.PlaylistID),
		FilterPlaylistID: filter.PlaylistID != 0,
		Limit:            int32(filter.Limit),
		Offset:           int32(filter.Offset),
	}

	tracks, err := t.repo.queries(ctx).TrackGetDeletedFilteredPopulated(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get filtered deleted tracks %+v | %w", filter, err)
	}

	return utils.SliceMap(tracks, func(t sqlc.TrackGetDeletedFilteredPopulatedRow) *model.Track {
		track := model.TrackModel(t.Track)
		track.Playlist = *model.PlaylistModel(t.Playlist)
		track.Playlist.Owner = *model.UserModel(t.User)
		track.DeletedAt = t.PlaylistTrack.DeletedAt.Time

		return track
	}), nil
}

func (t *Track) Create(ctx context.Context, track *model.Track) error {
	id, err := t.repo.queries(ctx).TrackCreate(ctx, sqlc.TrackCreateParams{
		SpotifyID:  track.SpotifyID,
		Name:       toString(track.Name),
		Popularity: toInt(track.Popularity),
		DurationMs: toInt(track.DurationMs),
	})
	if err != nil {
		return fmt.Errorf("create track %+v | %w", *track, err)
	}

	track.ID = int(id)

	return nil
}

func (t *Track) CreateArtist(ctx context.Context, artist *model.TrackArtist) error {
	id, err := t.repo.queries(ctx).TrackArtistCreate(ctx, sqlc.TrackArtistCreateParams{
		TrackID:  int32(artist.TrackID),
		ArtistID: int32(artist.ArtistID),
	})
	if err != nil {
		return fmt.Errorf("create track artist %+v | %w", *artist, err)
	}

	artist.ID = int(id)

	return nil
}

func (t *Track) Update(ctx context.Context, track model.Track) error {
	if err := t.repo.queries(ctx).TrackUpdate(ctx, sqlc.TrackUpdateParams{
		ID:         int32(track.ID),
		Name:       toString(track.Name),
		Popularity: toInt(track.Popularity),
		DurationMs: toInt(track.DurationMs),
	}); err != nil {
		return fmt.Errorf("update track %+v | %w", track, err)
	}

	return nil
}

func (t *Track) DeleteArtistByArtistTrack(ctx context.Context, artist model.TrackArtist) error {
	if err := t.repo.queries(ctx).TrackArtistDeleteByArtistTrack(ctx, sqlc.TrackArtistDeleteByArtistTrackParams{
		ArtistID: int32(artist.ArtistID),
		TrackID:  int32(artist.TrackID),
	}); err != nil {
		return fmt.Errorf("delete track artist %+v | %w", artist, err)
	}

	return nil
}
