package generator

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/internal/task"
	"github.com/topvennie/sortifyr/pkg/config"
)

const taskUID = "task-generator"

func getTaskUID(gen *model.Generator) string {
	if gen == nil {
		return taskUID
	}

	return taskUID + "-" + strconv.Itoa(gen.ID)
}

func getTaskName(gen *model.Generator) string {
	if gen == nil {
		return "Generator"
	}

	return "Generator - " + gen.Name
}

func (g *generator) taskRegister(ctx context.Context) error {
	if err := task.Manager.Add(ctx, task.NewTask(
		getTaskUID(nil),
		getTaskName(nil),
		config.GetDefaultDuration("task.generator_s", 60*60),
		func(ctx context.Context, users []model.User) []task.TaskResult {
			results := make([]task.TaskResult, 0, len(users))

			for _, user := range users {
				results = append(results, task.TaskResult{
					User:    user,
					Message: "",
					Error:   g.sync(ctx, user),
				})
			}

			return results
		},
	)); err != nil {
		return err
	}

	gens, err := g.generator.GetMaintainedPopulated(ctx)
	if err != nil {
		return err
	}

	for _, gen := range gens {
		if err := task.Manager.Add(ctx, task.NewTask(
			getTaskUID(gen),
			getTaskName(gen),
			gen.Interval,
			func(ctx context.Context, _ []model.User) []task.TaskResult {
				return []task.TaskResult{{
					User:    gen.User,
					Message: "",
					Error:   g.maintain(ctx, gen.User, gen.ID),
				}}
			},
		)); err != nil {
			return err
		}
	}

	return nil
}

// sync will update all the generator information
func (g *generator) sync(ctx context.Context, user model.User) error {
	gens, err := g.generator.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, gen := range gens {
		if err := g.syncOne(ctx, user, gen); err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) syncOne(ctx context.Context, user model.User, gen *model.Generator) error {
	gen.Outdated = false

	// If the generator has a playlist, does it still exist?
	var playlist *model.Playlist
	if gen.PlaylistID != 0 {
		playlists, err := g.playlist.GetByUser(ctx, user.ID)
		if err != nil {
			return err
		}
		idx := slices.IndexFunc(playlists, func(p *model.Playlist) bool { return p.ID == gen.PlaylistID })
		if idx == -1 {
			// Playlist is gone
			// The user probably deleted it manually
			gen.PlaylistID = 0
			gen.Maintained = false
			gen.Interval = 0
			gen.Outdated = false
		} else {
			playlist = playlists[idx]
		}
	}

	// If the generator still has a playlist, is it up to date?
	if playlist != nil {
		// Get the current tracks
		oldTracks, err := g.track.GetByPlaylist(ctx, playlist.ID)
		if err != nil {
			return err
		}

		// Get the new tracks
		newTracks, err := g.Generate(ctx, gen)
		if err != nil {
			return err
		}

		if len(oldTracks) != len(newTracks) {
			gen.Outdated = true
		} else {
			slices.SortFunc(oldTracks, func(a, b *model.Track) int { return a.ID - b.ID })
			slices.SortFunc(newTracks, func(a, b model.Track) int { return a.ID - b.ID })

			for i := range oldTracks {
				if !oldTracks[i].Equal(newTracks[i]) {
					gen.Outdated = true
					break
				}
			}
		}
	}

	if err := g.generator.Update(ctx, *gen); err != nil {
		return err
	}

	return nil
}

// maintain updates a maintained playlist to contain the newest tracks
func (g *generator) maintain(ctx context.Context, user model.User, genID int) error {
	// Refresh generator data
	gen, err := g.generator.Get(ctx, genID)
	if err != nil {
		return err
	}
	// Double check all values
	if !gen.Maintained {
		return errors.New("generator is no longer maintained")
	}

	playlist, err := g.playlist.Get(ctx, gen.PlaylistID)
	if err != nil {
		return err
	}
	if playlist == nil {
		return fmt.Errorf("no playlist with id %d", gen.PlaylistID)
	}

	// Get the current tracks
	oldTracks, err := g.track.GetByPlaylist(ctx, playlist.ID)
	if err != nil {
		return err
	}

	// Get the new tracks
	newTracks, err := g.Generate(ctx, gen)
	if err != nil {
		return err
	}

	// Find which ones to add and to remove
	toAdd := make([]model.Track, 0)
	toRemove := make([]model.Track, 0)

	for i := range newTracks {
		if idx := slices.IndexFunc(oldTracks, func(t *model.Track) bool { return newTracks[i].Equal(*t) }); idx == -1 {
			toAdd = append(toAdd, newTracks[i])
		}
	}

	for i := range oldTracks {
		if idx := slices.IndexFunc(newTracks, func(t model.Track) bool { return oldTracks[i].Equal(t) }); idx == -1 {
			toRemove = append(toRemove, *oldTracks[i])
		}
	}

	// Do the actions
	if err := spotifyapi.C.PlaylistDeleteTrackAll(ctx, user, playlist.SpotifyID, playlist.SnapshotID, toRemove); err != nil {
		return err
	}
	if err := spotifyapi.C.PlaylistPostTrackAll(ctx, user, playlist.SpotifyID, toAdd); err != nil {
		return err
	}

	gen.Outdated = false
	if err := g.generator.Update(ctx, *gen); err != nil {
		return err
	}

	return nil
}
