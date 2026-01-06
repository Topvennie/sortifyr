package generator

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/utils"
)

func (g *generator) Generate(ctx context.Context, gen *model.Generator) ([]model.Track, error) {
	var tracks []model.Track
	var err error

	normalize(gen)

	switch gen.Params.Preset {
	case model.GeneratorPresetTop:
		tracks, err = g.top(ctx, *gen)
	case model.GeneratorPresetOldTop:
		tracks, err = g.oldTop(ctx, *gen)

	default:
		tracks, err = g.custom(*gen)
	}

	if err != nil {
		return nil, err
	}

	return tracks, nil
}

type trackPlayCount struct {
	track     model.Track
	playCount int
}

func (g *generator) custom(gen model.Generator) ([]model.Track, error) {
	return nil, nil
}

func (g *generator) top(ctx context.Context, gen model.Generator) ([]model.Track, error) {
	params := gen.Params.ParamsTop

	// Get all excluded tracks
	// Can be from the excluded tracks list
	// Or the excluded playlists list
	excludedPlaylistTracks, err := g.playlist.GetTrackByPlaylistIDs(ctx, gen.Params.ExcludedPlaylistIDs)
	if err != nil {
		return nil, err
	}
	excludedTracksMap := make(map[int]bool)
	for _, t := range excludedPlaylistTracks {
		excludedTracksMap[t.TrackID] = true
	}
	for _, t := range gen.Params.ExcludedTrackIDs {
		excludedTracksMap[t] = true
	}

	// Get history for last 14 days
	skipped := false
	history, err := g.history.GetPopulatedFiltered(ctx, model.HistoryFilter{
		UserID:  gen.UserID,
		Start:   params.Window.Start,
		Skipped: &skipped,
	})
	if err != nil {
		return nil, err
	}

	playedAts := make(map[int][]time.Time)
	for _, h := range history {
		playedAt, ok := playedAts[h.TrackID]
		if !ok {
			playedAt = []time.Time{}
		}
		playedAt = append(playedAt, h.PlayedAt)
		playedAts[h.TrackID] = playedAt
	}

	tracks := make([]trackPlayCount, 0)
	seen := make(map[int]bool)
	for _, h := range history {
		if ok := seen[h.TrackID]; ok {
			continue
		}

		seen[h.TrackID] = true

		// Did the user exclude it
		if excludedTracksMap[h.TrackID] {
			continue
		}

		if hasBurst(playedAts[h.TrackID], params.Window) {
			tracks = append(tracks, trackPlayCount{
				track:     h.Track,
				playCount: len(playedAts[h.TrackID]),
			})
		}
	}

	slices.SortFunc(tracks, func(a, b trackPlayCount) int {
		if b.playCount == a.playCount {
			return strings.Compare(a.track.Name, b.track.Name)
		}

		return b.playCount - a.playCount
	})
	tracks = tracks[:min(gen.Params.TrackAmount, len(tracks))]
	slices.SortFunc(tracks, func(a, b trackPlayCount) int { return strings.Compare(a.track.Name, b.track.Name) })

	return utils.SliceMap(tracks, func(t trackPlayCount) model.Track { return t.track }), nil
}

func (g *generator) oldTop(ctx context.Context, gen model.Generator) ([]model.Track, error) {
	params := gen.Params.ParamsOldTop

	// Get all excluded tracks
	// Can be from the excluded tracks list
	// Or the excluded playlists list
	excludedPlaylistTracks, err := g.playlist.GetTrackByPlaylistIDs(ctx, gen.Params.ExcludedPlaylistIDs)
	if err != nil {
		return nil, err
	}
	excludedTracksMap := make(map[int]bool)
	for _, t := range excludedPlaylistTracks {
		excludedTracksMap[t.TrackID] = true
	}
	for _, t := range gen.Params.ExcludedTrackIDs {
		excludedTracksMap[t] = true
	}

	// Get the relevant recent history
	skipped := false
	recent, err := g.history.GetPopulatedFiltered(ctx, model.HistoryFilter{
		UserID:  gen.UserID,
		Start:   params.RecentWindow.Start,
		End:     params.RecentWindow.End,
		Skipped: &skipped,
	})
	if err != nil {
		return nil, err
	}

	// Get all play times for each track
	recentPlayedAts := make(map[int][]time.Time)
	for _, r := range recent {
		playedAt, ok := recentPlayedAts[r.TrackID]
		if !ok {
			playedAt = []time.Time{}
		}
		recentPlayedAts[r.TrackID] = append(playedAt, r.PlayedAt)
	}

	// Get all relevant peak history
	old, err := g.history.GetPopulatedFiltered(ctx, model.HistoryFilter{
		UserID:  gen.UserID,
		Start:   params.PeakWindow.Start,
		End:     params.PeakWindow.End,
		Skipped: &skipped,
	})
	if err != nil {
		return nil, err
	}

	// Get all play times for each track
	oldPlayedAts := make(map[int][]time.Time)
	for _, o := range old {
		playedAt, ok := oldPlayedAts[o.TrackID]
		if !ok {
			playedAt = []time.Time{}
		}
		oldPlayedAts[o.TrackID] = append(playedAt, o.PlayedAt)
	}

	tracks := make([]trackPlayCount, 0)
	seen := make(map[int]bool)
	for _, o := range old {
		if ok := seen[o.TrackID]; ok {
			continue
		}

		seen[o.TrackID] = true

		// Did the user exclude it
		if excludedTracksMap[o.TrackID] {
			continue
		}

		// Did we play it too much recently?
		if hasBurst(recentPlayedAts[o.TrackID], params.RecentWindow) {
			continue
		}

		// Did we play it enough times in the past?
		if !hasBurst(oldPlayedAts[o.TrackID], params.PeakWindow) {
			continue
		}

		tracks = append(tracks, trackPlayCount{
			track:     o.Track,
			playCount: len(oldPlayedAts[o.TrackID]),
		})
	}

	slices.SortFunc(tracks, func(a, b trackPlayCount) int {
		if b.playCount == a.playCount {
			return strings.Compare(a.track.Name, b.track.Name)
		}

		return b.playCount - a.playCount
	})
	tracks = tracks[:min(gen.Params.TrackAmount, len(tracks))]
	slices.SortFunc(tracks, func(a, b trackPlayCount) int { return strings.Compare(a.track.Name, b.track.Name) })

	return utils.SliceMap(tracks, func(t trackPlayCount) model.Track { return t.track }), nil
}
