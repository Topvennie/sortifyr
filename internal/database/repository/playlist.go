package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/pkg/sqlc"
	"github.com/topvennie/spotify_organizer/pkg/utils"
)

type Playlist struct {
	repo Repository
}

func (r *Repository) NewPlaylist() *Playlist {
	return &Playlist{
		repo: *r,
	}
}

func (p *Playlist) GetAllPopulated(ctx context.Context) ([]*model.Playlist, error) {
	playlists, err := p.repo.queries(ctx).PlaylistGetAllWithOwner(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all playlists %w", err)
	}

	return utils.SliceMap(playlists, func(r sqlc.PlaylistGetAllWithOwnerRow) *model.Playlist {
		return model.PlaylistModelPopulated(r.Playlist, r.User)
	}), nil
}

func (p *Playlist) Create(ctx context.Context, playlist *model.Playlist) error {
	id, err := p.repo.queries(ctx).PlaylistCreate(ctx, sqlc.PlaylistCreateParams{
		SpotifyID:     playlist.SpotifyID,
		OwnerID:       int32(playlist.OwnerID),
		Name:          playlist.Name,
		Description:   pgtype.Text{String: playlist.Description, Valid: playlist.Description != ""},
		Public:        playlist.Public,
		Tracks:        int32(playlist.Tracks),
		Collaborative: playlist.Collaborative,
	})
	if err != nil {
		return fmt.Errorf("create playlist %+v | %w", *playlist, err)
	}

	playlist.ID = int(id)

	return nil
}

func (p *Playlist) Update(ctx context.Context, playlist model.Playlist) error {
	if err := p.repo.queries(ctx).PlaylistUpdate(ctx, sqlc.PlaylistUpdateParams{
		ID:            int32(playlist.ID),
		OwnerID:       int32(playlist.OwnerID),
		Name:          playlist.Name,
		Description:   pgtype.Text{String: playlist.Description, Valid: playlist.Description != ""},
		Public:        playlist.Public,
		Tracks:        int32(playlist.Tracks),
		Collaborative: playlist.Collaborative,
	}); err != nil {
		return fmt.Errorf("update playlist %+v | %w", playlist, err)
	}

	return nil
}

func (p *Playlist) Delete(ctx context.Context, id int) error {
	if err := p.repo.queries(ctx).PlaylistDelete(ctx, int32(id)); err != nil {
		return fmt.Errorf("delete playlist %d", id)
	}

	return nil
}
