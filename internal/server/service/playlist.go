package service

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
	"github.com/topvennie/sortifyr/internal/server/dto"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/internal/task"
	"github.com/topvennie/sortifyr/pkg/storage"
	"github.com/topvennie/sortifyr/pkg/utils"
	"go.uber.org/zap"
)

const taskPlaylistDuplicateUID = "task-playlist-duplicate"

type Playlist struct {
	service Service

	playlist repository.Playlist
	user     repository.User
}

func (s *Service) NewPlaylist() *Playlist {
	return &Playlist{
		service:  *s,
		playlist: *s.repo.NewPlaylist(),
		user:     *s.repo.NewUser(),
	}
}

func (p *Playlist) GetByUser(ctx context.Context, userID int) ([]dto.Playlist, error) {
	playlistsDB, err := p.playlist.GetByUserPopulated(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if playlistsDB == nil {
		return []dto.Playlist{}, nil
	}

	return utils.SliceMap(playlistsDB, func(p *model.Playlist) dto.Playlist { return dto.PlaylistDTO(p, &p.Owner) }), nil
}

func (p *Playlist) GetCover(ctx context.Context, playlistID int) ([]byte, error) {
	playlist, err := p.playlist.Get(ctx, playlistID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if playlist == nil {
		return nil, fiber.ErrNotFound
	}
	if playlist.CoverID == "" {
		return nil, fiber.ErrNotFound
	}

	cover, err := storage.S.Get(playlist.CoverID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return cover, nil
}

func (p *Playlist) GetDuplicates(ctx context.Context, userID int) ([]dto.PlaylistDuplicate, error) {
	playlistsAll, err := p.playlist.GetDuplicateTracksByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	// Filter out tracks without spotify id
	// We can't delete them with the api
	playlists := make([]*model.Playlist, 0)
	for i := range playlistsAll {
		playlistsAll[i].Duplicates = utils.SliceFilter(playlistsAll[i].Duplicates, func(t model.Track) bool { return t.SpotifyID != "" })
		if len(playlistsAll[i].Duplicates) > 0 {
			playlists = append(playlists, playlistsAll[i])
		}
	}

	return utils.SliceMap(playlists, func(p *model.Playlist) dto.PlaylistDuplicate {
		return dto.PlaylistDuplicateDTO(p, &p.Owner, p.Duplicates)
	}), nil
}

func (p *Playlist) GetUnplayables(ctx context.Context, userID int) ([]dto.PlaylistUnplayable, error) {
	playlists, err := p.playlist.GetUnplayableTracksByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(playlists, func(p *model.Playlist) dto.PlaylistUnplayable {
		return dto.PlaylistUnplayableDTO(p, &p.Owner, p.Unplayables)
	}), nil
}

func (p *Playlist) RemoveDuplicates(ctx context.Context, userID int) error {
	user, err := p.user.GetByID(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return fiber.ErrUnauthorized
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		taskPlaylistDuplicateUID,
		"Playlist Duplicates Removal",
		task.IntervalOnce,
		func(ctx context.Context, _ []model.User) []task.TaskResult {
			return []task.TaskResult{{
				User:    *user,
				Message: "",
				Error:   p.removeDuplicatesTask(ctx, *user),
			}}
		},
	)); err != nil {
		if errors.Is(err, task.ErrTaskExists) {
			return fiber.NewError(fiber.StatusBadRequest, "Task is already running")
		}
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (p *Playlist) removeDuplicatesTask(ctx context.Context, user model.User) error {
	// Spotify's API will remove every instance of a track in a playlist
	// That's why we first need to make an API call to delete them
	// Which deleted all instances
	// And then add them back so that we dont lose them.
	playlists, err := p.playlist.GetDuplicateTracksByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for i := range playlists {
		if playlists[i].SnapshotID == "" {
			continue
		}

		unique := utils.SliceUniqueFunc(playlists[i].Duplicates, func(t model.Track) string { return t.SpotifyID })
		filtered := utils.SliceFilter(unique, func(t model.Track) bool { return t.SpotifyID != "" })

		// Remove them all
		if err := spotifyapi.C.PlaylistDeleteTrackAll(ctx, user, playlists[i].SpotifyID, playlists[i].SnapshotID, filtered); err != nil {
			return err
		}

		// Add them all back
		if err := spotifyapi.C.PlaylistPostTrackAll(ctx, user, playlists[i].SpotifyID, filtered); err != nil {
			return err
		}
	}

	return nil
}
