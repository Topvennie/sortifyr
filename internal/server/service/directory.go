package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/internal/database/repository"
	"github.com/topvennie/spotify_organizer/internal/server/dto"
	"github.com/topvennie/spotify_organizer/pkg/utils"
	"go.uber.org/zap"
)

type Directory struct {
	service Service

	directory repository.Directory
}

func (s *Service) NewDirectory() *Directory {
	return &Directory{
		service:   *s,
		directory: *s.repo.NewDirectory(),
	}
}

func (d *Directory) GetByUser(ctx context.Context, userID int) ([]dto.Directory, error) {
	directoryModels, err := d.directory.GetByUserPopulated(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if directoryModels == nil {
		return []dto.Directory{}, nil
	}

	roots := utils.SliceFilter(directoryModels, func(d *model.Directory) bool { return d.ParentID == 0 })
	directories := make([]dto.Directory, 0, len(roots))

	for _, root := range roots {
		directories = append(directories, dto.DirectoryDTO(root, directoryModels))
	}

	return directories, nil
}

// Sync brings the database up to date with the data received from the api
// nolint:gocognit // It's fine
func (d *Directory) Sync(ctx context.Context, userID int, roots []dto.Directory) ([]dto.Directory, error) {
	directoryDTOs := make([]dto.Directory, 0)

	now := make([]dto.Directory, 0, len(roots))
	next := make([]dto.Directory, 0)

	now = append(now, roots...)

	for len(now) > 0 {
		for _, n := range now {
			directoryDTOs = append(directoryDTOs, n)
			if len(n.Children) > 0 {
				next = append(next, n.Children...)
			}
		}

		now = next
		next = []dto.Directory{}
	}

	directories := utils.SliceMap(directoryDTOs, func(d dto.Directory) *model.Directory { return d.ToModel(userID, directoryDTOs) })
	directoriesDB, err := d.directory.GetByUserPopulated(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	toCreate := make([]*model.Directory, 0)
	toUpdate := make([]struct {
		new *model.Directory
		old *model.Directory
	}, 0)
	toDelete := make([]*model.Directory, 0)

	// Let's get all directories that need to be created or updated
	for _, directory := range directories {
		// If the directory doesn't have an id yet then it needs to be created
		if directory.ID == 0 {
			toCreate = append(toCreate, directory)
			continue
		}

		// It might be an update, let's check the values
		directoryDB, ok := utils.SliceFind(directoriesDB, func(d *model.Directory) bool { return d.ID == directory.ID })
		zap.S().Debug(*directoryDB)
		if !ok {
			// User gave an invalid id
			// Unlucky for them
			continue
		}

		if !(*directoryDB).Equal(*directory) {
			// Not equal so let's add it to the update list
			toUpdate = append(toUpdate, struct {
				new *model.Directory
				old *model.Directory
			}{
				new: directory,
				old: *directoryDB,
			})
		}
	}

	// Get all directories that need to be deleted
	for _, directorDB := range directoriesDB {
		if _, ok := utils.SliceFind(directories, func(d *model.Directory) bool { return d.ID == directorDB.ID }); !ok {
			toDelete = append(toDelete, directorDB)
		}
	}

	// All the directories are sorted
	// Time to do the database operations
	if err := d.service.withRollback(ctx, func(ctx context.Context) error {
		// Create
		for _, directory := range toCreate {
			if err := d.directory.Create(ctx, directory); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}
		}

		// Update
		for _, entry := range toUpdate {
			if err := d.directory.Update(ctx, *entry.new); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}

			// Create / Delete linked playlists
			for i := range entry.new.Playlists {
				if _, ok := utils.SliceFind(entry.old.Playlists, func(p model.Playlist) bool { return p.ID == entry.new.Playlists[i].ID }); !ok {
					if err := d.directory.CreatePlaylist(ctx, &model.DirectoryPlaylist{
						DirectoryID: entry.new.ID,
						PlaylistID:  entry.new.Playlists[i].ID,
					}); err != nil {
						zap.S().Error(err)
						return fiber.ErrInternalServerError
					}
				}
			}

			for i := range entry.old.Playlists {
				if _, ok := utils.SliceFind(entry.new.Playlists, func(p model.Playlist) bool { return p.ID == entry.old.Playlists[i].ID }); !ok {
					if err := d.directory.DeletePlaylist(ctx, entry.old.Playlists[i].ID); err != nil {
						zap.S().Error(err)
						return fiber.ErrInternalServerError
					}
				}
			}
		}

		// Delete
		for _, directory := range toDelete {
			if err := d.directory.Delete(ctx, directory.ID); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return d.GetByUser(ctx, userID)
}
