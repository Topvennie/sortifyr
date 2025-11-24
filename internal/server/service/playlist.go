package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/internal/database/repository"
	"github.com/topvennie/spotify_organizer/internal/server/dto"
	"github.com/topvennie/spotify_organizer/internal/spotify"
	"github.com/topvennie/spotify_organizer/pkg/utils"
	"go.uber.org/zap"
)

type Playlist struct {
	service Service

	playlist repository.Playlist
}

func (s *Service) NewPlaylist() *Playlist {
	return &Playlist{
		service:  *s,
		playlist: *s.repo.NewPlaylist(),
	}
}

func (p *Playlist) GetAll(ctx context.Context) ([]dto.Playlist, error) {
	playlistsDB, err := p.playlist.GetAllPopulated(ctx)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if playlistsDB == nil {
		return []dto.Playlist{}, nil
	}

	return utils.SliceMap(playlistsDB, func(p *model.Playlist) dto.Playlist { return dto.PlaylistDTO(p, &p.Owner) }), nil
}

func (p *Playlist) Sync(ctx context.Context, uid string) error {
	if err := spotify.C.PlaylistSync(ctx, uid); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
