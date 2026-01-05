package generator

import (
	"context"
	"errors"
	"slices"
	"strconv"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/internal/spotifysync"
	"github.com/topvennie/sortifyr/internal/task"
	"github.com/topvennie/sortifyr/pkg/config"
	"github.com/topvennie/sortifyr/pkg/utils"
	"go.uber.org/zap"
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
		config.GetDefaultDurationS("task.generator_s", 60*60),
		false,
		func(ctx context.Context, users []model.User) []task.TaskResult {
			results := make([]task.TaskResult, 0, len(users))

			for _, user := range users {
				results = append(results, task.TaskResult{
					User:    user,
					Message: "",
					Error:   g.spotifyStatus(ctx, user),
				})
			}

			return results
		},
	)); err != nil {
		return err
	}

	gens, err := g.generator.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, gen := range gens {
		if gen.Interval == 0 {
			continue
		}

		if err := task.Manager.Add(ctx, task.NewTask(
			getTaskUID(gen),
			getTaskName(gen),
			gen.Interval,
			true,
			func(ctx context.Context, _ []model.User) []task.TaskResult {
				return []task.TaskResult{{
					User:    gen.User,
					Message: "",
					Error:   g.refresh(ctx, gen.User, gen.ID),
				}}
			},
		)); err != nil {
			return err
		}
	}

	return nil
}

// spotifyStatus will check if the spotify playlist is up to date
// It does NOT update the tracks from the generator
func (g *generator) spotifyStatus(ctx context.Context, user model.User) error {
	gens, err := g.generator.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, gen := range gens {
		if err := g.spotifyStatusOne(ctx, user, gen); err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) spotifyStatusOne(ctx context.Context, user model.User, gen *model.Generator) error {
	if gen.PlaylistID == 0 {
		// Generator doesn't have a spotify playlist
		return nil
	}

	gen.SpotifyOutdated = false

	var playlist *model.Playlist
	// Get all the user's playlists
	playlists, err := g.playlist.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}
	// Find the playlist
	idx := slices.IndexFunc(playlists, func(p *model.Playlist) bool { return p.ID == gen.PlaylistID })
	if idx == -1 {
		// Playlist is gone
		// The user probably deleted it manually
		gen.PlaylistID = 0
		gen.SpotifyOutdated = true
		zap.S().Debug("No playlist")
	} else {
		playlist = playlists[idx]
	}

	// Is the playlist still up to date?
	if playlist != nil {
		// Get the playlist tracks
		playlistTracks, err := g.track.GetByPlaylist(ctx, playlist.ID)
		if err != nil {
			return err
		}

		slices.SortFunc(playlistTracks, func(a, b *model.Track) int { return a.ID - b.ID })
		slices.SortFunc(gen.Tracks, func(a, b model.Track) int { return a.ID - b.ID })

		gen.SpotifyOutdated = !slices.EqualFunc(playlistTracks, gen.Tracks, func(a *model.Track, b model.Track) bool { return a.Equal(b) })
		zap.S().Debug("slices euqla")
		zap.S().Debug(!slices.EqualFunc(playlistTracks, gen.Tracks, func(a *model.Track, b model.Track) bool { return a.Equal(b) }))
	}

	if err := g.generator.Update(ctx, *gen); err != nil {
		return err
	}

	return nil
}

// refresh will refresh the generator tracks and update the Spotify playlist if applicable
func (g *generator) refresh(ctx context.Context, user model.User, genID int) error {
	gen, err := g.generator.Get(ctx, genID)
	if err != nil {
		return err
	}

	newTracks, err := G.Generate(ctx, gen)
	if err != nil {
		return err
	}
	slices.SortFunc(newTracks, func(a, b model.Track) int { return a.ID - b.ID })

	// Update the generator database tracks
	dbTracks, err := g.track.GetByGenerator(ctx, gen.ID)
	if err != nil {
		return err
	}
	slices.SortFunc(dbTracks, func(a, b *model.Track) int { return a.ID - b.ID })

	// Update the db if needed
	if equal := slices.EqualFunc(newTracks, dbTracks, func(a model.Track, b *model.Track) bool { return b.Equal(a) }); !equal {
		if err := g.generator.DeleteTrackByGenerator(ctx, gen.ID); err != nil {
			return err
		}
		if err := g.generator.CreateTrackBatch(ctx, utils.SliceMap(newTracks, func(t model.Track) model.GeneratorTrack {
			return model.GeneratorTrack{GeneratorID: gen.ID, TrackID: t.ID}
		})); err != nil {
			return err
		}
	}

	// Update the Spotify playlist
	if gen.PlaylistID != 0 {
		playlist, err := g.playlist.Get(ctx, gen.PlaylistID)
		if err != nil {
			return err
		}
		if playlist == nil {
			return errors.New("db unsyned")
		}

		// Get the latest playlist tracks
		playlistTracksAPI, err := spotifyapi.C.PlaylistGetTrackAll(ctx, user, playlist.SpotifyID)
		if err != nil {
			return err
		}
		playlistTracks := utils.SliceMap(playlistTracksAPI, func(t spotifyapi.Track) model.Track { return t.ToModel() })

		toCreate := []model.Track{}
		toDelete := []model.Track{}
		for i := range newTracks {
			if idx := slices.IndexFunc(playlistTracks, func(t model.Track) bool { return t.Equal(newTracks[i]) }); idx == -1 {
				toCreate = append(toCreate, newTracks[i])
			}
		}
		for i := range playlistTracks {
			if idx := slices.IndexFunc(newTracks, func(t model.Track) bool { return playlistTracks[i].Equal(t) }); idx == -1 {
				toDelete = append(toDelete, playlistTracks[i])
			}
		}

		if err := spotifyapi.C.PlaylistDeleteTrackAll(ctx, user, playlist.SpotifyID, playlist.SnapshotID, toDelete); err != nil {
			return err
		}
		if err := spotifyapi.C.PlaylistPostTrackAll(ctx, user, playlist.SpotifyID, toCreate); err != nil {
			return err
		}

		if len(toCreate) > 0 || len(toDelete) > 0 {
			if err := task.Manager.RunRecurringByUID(spotifysync.TaskPlaylistUID, user); err != nil {
				return err
			}
		}
	}

	gen.SpotifyOutdated = false
	if err := g.generator.Update(ctx, *gen); err != nil {
		return err
	}

	return nil
}
