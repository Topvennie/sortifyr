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

func (g *generator) Create(ctx context.Context, gen *model.Generator, createPlaylist bool) error {
	// Get user
	user, err := g.user.GetByID(ctx, gen.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found %d", gen.UserID)
	}

	// Create playlist
	gen.PlaylistID = 0
	if createPlaylist {
		// Get tracks
		tracks, err := g.Generate(ctx, gen)
		if err != nil {
			return err
		}

		// Create playlist in Spotify
		description := "Created by Sortifyr"
		if gen.Maintained {
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
		if err := g.playlist.Create(ctx, &playlist); err != nil {
			return err
		}
		if err := g.playlist.CreateUser(ctx, &model.PlaylistUser{UserID: user.ID, PlaylistID: playlist.ID}); err != nil {
			return err
		}

		gen.PlaylistID = playlist.ID

		// Add tracks to the playlist
		if err := spotifyapi.C.PlaylistPostTrackAll(ctx, *user, playlist.SpotifyID, tracks); err != nil {
			return err
		}
	} else {
		normalize(gen)
	}

	// Save in database
	if err := g.generator.Create(ctx, gen); err != nil {
		return err
	}

	// Start task for maintaince
	if gen.Maintained {
		if gen.Interval.Nanoseconds() == 0 {
			return fmt.Errorf("interval is equal to 0 %+v", *gen)
		}

		if err := task.Manager.Add(ctx, task.NewTask(
			getTaskUID(gen),
			getTaskName(gen),
			gen.Interval,
			true,
			func(ctx context.Context, _ []model.User) []task.TaskResult {
				return []task.TaskResult{{
					User:    *user,
					Message: "",
					Error:   g.maintain(ctx, *user, gen.ID),
				}}
			},
		)); err != nil {
			return err
		}
	}

	return nil
}

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

	// We're going to be very lazy
	// If the old generator had a playlist, delete it
	// If the new generator has a playlist, create it
	// This will sometimes delete and create a playlist after one another

	// Delete the old playlist
	if oldGen.PlaylistID != 0 {
		playlist, err := g.playlist.Get(ctx, oldGen.PlaylistID)
		if err != nil {
			return err
		}
		if playlist == nil {
			return fmt.Errorf("playlist with id %d not found", oldGen.PlaylistID)
		}

		if err := g.playlist.DeleteUserByUserPlaylist(ctx, model.PlaylistUser{UserID: user.ID, PlaylistID: playlist.ID}); err != nil {
			return err
		}
		if err := spotifyapi.C.PlaylistDelete(ctx, *user, playlist.SpotifyID); err != nil {
			return err
		}
	}

	// Create the new playlist
	gen.PlaylistID = 0
	if createPlaylist {
		// Get tracks
		tracks, err := g.Generate(ctx, gen)
		if err != nil {
			return err
		}

		// Create playlist in Spotify
		description := "Created by Sortifyr"
		if gen.Maintained {
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

		gen.PlaylistID = playlist.ID

		// Add tracks to the playlist
		if err := spotifyapi.C.PlaylistPostTrackAll(ctx, *user, playlist.SpotifyID, tracks); err != nil {
			return err
		}
	} else {
		normalize(gen)
	}

	// Update in database
	if err := g.generator.Update(ctx, *gen); err != nil {
		return err
	}

	// Remove old task
	if oldGen.Maintained {
		if err := task.Manager.Remove(ctx, getTaskUID(oldGen)); err != nil {
			if !errors.Is(err, task.ErrTaskNotExists) {
				return err
			}
		}
	}

	// Start task for maintaince
	if gen.Maintained {
		if gen.Interval.Nanoseconds() == 0 {
			return fmt.Errorf("interval is equal to 0 %+v", *gen)
		}

		if err := task.Manager.Add(ctx, task.NewTask(
			getTaskUID(gen),
			getTaskName(gen),
			gen.Interval,
			true,
			func(ctx context.Context, _ []model.User) []task.TaskResult {
				return []task.TaskResult{{
					User:    *user,
					Message: "",
					Error:   g.maintain(ctx, *user, gen.ID),
				}}
			},
		)); err != nil {
			return err
		}
	}

	return nil
}
