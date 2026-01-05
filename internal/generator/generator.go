// Package generator creates and maintaines playlists based on parameters and presets
package generator

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/internal/task"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type generator struct {
	generator repository.Generator
	history   repository.History
	playlist  repository.Playlist
	track     repository.Track
	user      repository.User
}

var G *generator

func Init(repo repository.Repository) error {
	G = &generator{
		generator: *repo.NewGenerator(),
		history:   *repo.NewHistory(),
		playlist:  *repo.NewPlaylist(),
		track:     *repo.NewTrack(),
		user:      *repo.NewUser(),
	}

	if err := G.taskRegister(context.Background()); err != nil {
		return err
	}

	return nil
}

func (g *generator) Refresh(ctx context.Context, user model.User, gen model.Generator) error {
	// If the generator has an interval then it has a scheduled task to update it.
	// So we can just run that.
	if gen.Interval > 0 {
		return task.Manager.RunRecurringByUID(getTaskUID(&gen), user)
	}

	// Else we need to add a one time task to
	return task.Manager.Add(ctx, task.NewTask(
		getTaskUID(&gen),
		getTaskName(&gen),
		task.IntervalOnce,
		true,
		func(ctx context.Context, _ []model.User) []task.TaskResult {
			return []task.TaskResult{{
				User:    user,
				Message: "",
				Error:   g.refresh(ctx, user, gen.ID),
			}}
		},
	))
}

func (g *generator) Create(ctx context.Context, gen *model.Generator, createPlaylist bool) error {
	// Get user
	user, err := g.user.GetByID(ctx, gen.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found %d", gen.UserID)
	}

	// Create playlist in spotify and the db
	gen.PlaylistID = 0
	if createPlaylist {
		// Create playlist in Spotify
		description := "Created by Sortifyr"
		if gen.Interval > 0 {
			days := int(gen.Interval.Nanoseconds() / int64(24*time.Hour))
			daysStr := strconv.Itoa(days) + "days"
			if days == 1 {
				daysStr = "day"
			}
			description = fmt.Sprintf("Created and maintained (every %s) by Sortifyr", daysStr)
		}
		public := false
		collaborative := false

		playlist := model.Playlist{
			OwnerID:       gen.UserID,
			Name:          gen.Name,
			Description:   description,
			Public:        &public,
			Collaborative: &collaborative,
		}

		if err := spotifyapi.C.PlaylistCreate(ctx, *user, &playlist); err != nil {
			return err
		}

		// Create playlist in db and link to user
		if err := g.playlist.Create(ctx, &playlist); err != nil {
			return err
		}
		if err := g.playlist.CreateUser(ctx, &model.PlaylistUser{UserID: user.ID, PlaylistID: playlist.ID}); err != nil {
			return err
		}

		gen.PlaylistID = playlist.ID
	}

	// Save in database
	if err := g.generator.Create(ctx, gen); err != nil {
		return err
	}

	// Add tracks to the playlist
	// If the generator has an interval then the task will do it on the first run
	// Else we need to just run the task once
	interval := task.IntervalOnce
	if gen.Interval > 0 {
		interval = gen.Interval
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		getTaskUID(gen),
		getTaskName(gen),
		interval,
		true,
		func(ctx context.Context, _ []model.User) []task.TaskResult {
			return []task.TaskResult{{
				User:    *user,
				Message: "",
				Error:   g.refresh(ctx, *user, gen.ID),
			}}
		},
	)); err != nil {
		return err
	}

	return nil
}

// nolint:gocognit // It's fine
func (g *generator) Update(ctx context.Context, gen *model.Generator, createPlaylist bool) error {
	// Get user
	user, err := g.user.GetByID(ctx, gen.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found %d", gen.UserID)
	}

	// Get old generator
	oldGen, err := g.generator.Get(ctx, gen.ID)
	if err != nil {
		return err
	}
	if oldGen == nil {
		return fmt.Errorf("gen with id %d not found", gen.ID)
	}

	if oldGen.PlaylistID != 0 {
		playlist, err := g.playlist.Get(ctx, oldGen.PlaylistID)
		if err != nil {
			return err
		}
		if playlist == nil {
			return fmt.Errorf("playlist with id %d not found", oldGen.PlaylistID)
		}

		if createPlaylist {
			// Delete all tracks
			tracks, err := g.track.GetByPlaylist(ctx, playlist.ID)
			if err != nil {
				return err
			}

			if err := spotifyapi.C.PlaylistDeleteTrackAll(ctx, *user, playlist.SpotifyID, playlist.SnapshotID, utils.SliceDereference(tracks)); err != nil {
				return err
			}
		} else {
			// Delete playlist
			// Delete it from our db
			if err := g.playlist.DeleteUserByUserPlaylist(ctx, model.PlaylistUser{UserID: user.ID, PlaylistID: playlist.ID}); err != nil {
				return err
			}
			// Delete it from Spotify
			if err := spotifyapi.C.PlaylistDelete(ctx, *user, playlist.SpotifyID); err != nil {
				return err
			}
		}
	}

	// Create the new playlist
	gen.PlaylistID = 0
	if createPlaylist {
		if oldGen.PlaylistID == 0 {
			// Create playlist in Spotify
			description := "Created by Sortifyr"
			if gen.Interval > 0 {
				days := int(gen.Interval.Nanoseconds() / int64(24*time.Hour))
				daysStr := strconv.Itoa(days) + "days"
				if days == 1 {
					daysStr = "day"
				}
				description = fmt.Sprintf("Created and maintained (every %s) by Sortifyr", daysStr)
			}
			public := false
			collaborative := false

			playlist := model.Playlist{
				OwnerID:       gen.UserID,
				Name:          gen.Name,
				Description:   description,
				Public:        &public,
				Collaborative: &collaborative,
			}

			if err := spotifyapi.C.PlaylistCreate(ctx, *user, &playlist); err != nil {
				return err
			}

			// Create playlist in db and link to user
			if err := g.playlist.Create(ctx, &playlist); err != nil {
				return err
			}
			if err := g.playlist.CreateUser(ctx, &model.PlaylistUser{UserID: user.ID, PlaylistID: playlist.ID}); err != nil {
				return err
			}

			gen.PlaylistID = playlist.ID
		} else {
			// Use the old playlist
			gen.PlaylistID = oldGen.PlaylistID
		}
	}

	// Update in database
	if err := g.generator.Update(ctx, *gen); err != nil {
		return err
	}

	// Remove old task
	if oldGen.Interval > 0 {
		if err := task.Manager.Remove(ctx, getTaskUID(oldGen)); err != nil {
			if !errors.Is(err, task.ErrTaskNotExists) {
				return err
			}
		}
	}

	// Add tracks to the playlist
	// If the generator is maintained then the task will do it on the first run
	// Else we need to just run the task once
	interval := task.IntervalOnce
	if gen.Interval > 0 {
		interval = gen.Interval
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		getTaskUID(gen),
		getTaskName(gen),
		interval,
		true,
		func(ctx context.Context, _ []model.User) []task.TaskResult {
			return []task.TaskResult{{
				User:    *user,
				Message: "",
				Error:   g.refresh(ctx, *user, gen.ID),
			}}
		},
	)); err != nil {
		return err
	}

	return nil
}

func (g *generator) Delete(ctx context.Context, user model.User, gen model.Generator, deletePlaylist bool) error {
	if gen.PlaylistID != 0 && deletePlaylist {
		playlist, err := g.playlist.Get(ctx, gen.PlaylistID)
		if err != nil {
			return err
		}
		if playlist == nil {
			return fmt.Errorf("db unsynced %+v | %w", gen, err)
		}

		if err := spotifyapi.C.PlaylistDelete(ctx, user, playlist.SpotifyID); err != nil {
			return fmt.Errorf("delete playlist %w", err)
		}

		if err := g.playlist.DeleteUserByUserPlaylist(ctx, model.PlaylistUser{UserID: user.ID, PlaylistID: playlist.ID}); err != nil {
			return err
		}
	}

	if err := g.generator.Delete(ctx, gen.ID); err != nil {
		return err
	}

	return nil
}
