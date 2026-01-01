// Package generator creates and maintaines playlists based on parameters and presets
package generator

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type generator struct {
	history repository.History
}

var G *generator

func Init(repo repository.Repository) {
	G = &generator{
		history: *repo.NewHistory(),
	}
}

func (g *generator) Generate(ctx context.Context, preset model.GeneratorPreset, params model.GeneratorParameters) ([]model.Track, error) {
	var tracks []model.Track
	var err error

	switch preset {
	case model.PresetForgotten:
		tracks, err = g.forgotten(params)
	case model.PresetTop:
		tracks, err = g.top(ctx, params)
	case model.PresetOldTop:
		tracks, err = g.oldTop(ctx, params)

	default:
		tracks, err = g.custom(params)
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

func (g *generator) custom(params model.GeneratorParameters) ([]model.Track, error) {
	return nil, nil
}

func (g *generator) forgotten(params model.GeneratorParameters) ([]model.Track, error) {
	return nil, nil
}

func (g *generator) top(ctx context.Context, params model.GeneratorParameters) ([]model.Track, error) {
	// Set default param values
	if params.TrackAmount == 0 {
		params.TrackAmount = 50
	}
	if params.MinPlayCount == 0 {
		params.MinPlayCount = 5
	}

	// Get history for last 14 days
	skipped := false
	history, err := g.history.GetPopulatedFiltered(ctx, model.HistoryFilter{
		UserID:  params.UserID,
		Start:   time.Now().Add(-14 * 24 * time.Hour),
		Skipped: &skipped,
	})
	if err != nil {
		return nil, err
	}

	trackMap := make(map[int]trackPlayCount)
	for _, h := range history {
		track, ok := trackMap[h.TrackID]
		if !ok {
			track = trackPlayCount{track: h.Track, playCount: 0}
		}
		track.playCount++
		trackMap[h.TrackID] = track
	}

	tracks := utils.MapValues(trackMap)

	tracks = utils.SliceFilter(tracks, func(t trackPlayCount) bool { return t.playCount > params.MinPlayCount })
	slices.SortFunc(tracks, func(a, b trackPlayCount) int {
		if b.playCount == a.playCount {
			return strings.Compare(a.track.Name, b.track.Name)
		}

		return b.playCount - a.playCount
	})
	tracks = tracks[:min(params.TrackAmount, len(tracks))]
	slices.SortFunc(tracks, func(a, b trackPlayCount) int { return strings.Compare(a.track.Name, b.track.Name) })

	return utils.SliceMap(tracks, func(t trackPlayCount) model.Track { return t.track }), nil
}

func (g *generator) oldTop(ctx context.Context, params model.GeneratorParameters) ([]model.Track, error) {
	// Set default param values
	// This function assumes that the peak and recent windows fields are either fully populated or not at all
	now := time.Now()
	if params.PeakWindow.Start.IsZero() {
		params.PeakWindow.Start = now.Add(-1 * 24 * 365 * time.Hour) // 365 days ago
		params.PeakWindow.End = now.Add(-1 * 24 * 100 * time.Hour)   // 100 days ago
		params.PeakWindow.MinPlays = 5
		params.PeakWindow.BurstInterval = 24 * 14 * time.Hour // 14 days
	}
	if params.RecentWindow.Start.IsZero() {
		params.RecentWindow.Start = now.Add(-1 * 24 * 14 * time.Hour) // 14 days ago
		params.RecentWindow.End = now                                 // now
		params.RecentWindow.MinPlays = 2
		params.RecentWindow.BurstInterval = 24 * 14 * time.Hour // 14 days
	}

	skipped := false

	// Get the relevant recent history
	recent, err := g.history.GetPopulatedFiltered(ctx, model.HistoryFilter{
		UserID:  params.UserID,
		Start:   params.RecentWindow.Start,
		End:     params.RecentWindow.End,
		Skipped: &skipped,
	})
	if err != nil {
		return nil, err
	}

	// Get all play times for each track
	recentPlayedAt := make(map[int][]time.Time)
	for _, r := range recent {
		playedAt, ok := recentPlayedAt[r.TrackID]
		if !ok {
			playedAt = []time.Time{}
		}
		recentPlayedAt[r.TrackID] = append(playedAt, r.PlayedAt)
	}

	// Get all relevant peak history
	old, err := g.history.GetPopulatedFiltered(ctx, model.HistoryFilter{
		UserID:  params.UserID,
		Start:   params.PeakWindow.Start,
		End:     params.PeakWindow.End,
		Skipped: &skipped,
	})
	if err != nil {
		return nil, err
	}

	// Get a map from track id to track
	oldMap := make(map[int]model.Track)
	for _, o := range old {
		if _, ok := oldMap[o.TrackID]; !ok {
			oldMap[o.TrackID] = o.Track
		}
	}

	// Get all play times for each track
	oldPlayedAt := make(map[int][]time.Time)
	for _, o := range old {
		playedAt, ok := oldPlayedAt[o.TrackID]
		if !ok {
			playedAt = []time.Time{}
		}
		oldPlayedAt[o.TrackID] = append(playedAt, o.PlayedAt)
	}

	trackMap := make(map[int]trackPlayCount)
	for k, v := range oldPlayedAt {
		// Did we play it too much recently?
		if hasBurst(recentPlayedAt[k], params.RecentWindow.MinPlays, params.RecentWindow.BurstInterval) {
			continue
		}

		// Did we play it enough times in the past?
		if !hasBurst(v, params.PeakWindow.MinPlays, params.PeakWindow.BurstInterval) {
			continue
		}

		trackMap[k] = trackPlayCount{
			track:     oldMap[k],
			playCount: len(oldPlayedAt[k]),
		}
	}

	tracks := utils.MapValues(trackMap)
	slices.SortFunc(tracks, func(a, b trackPlayCount) int {
		if b.playCount == a.playCount {
			return strings.Compare(a.track.Name, b.track.Name)
		}

		return b.playCount - a.playCount
	})
	tracks = tracks[:min(params.TrackAmount, len(tracks))]
	slices.SortFunc(tracks, func(a, b trackPlayCount) int { return strings.Compare(a.track.Name, b.track.Name) })

	return utils.SliceMap(tracks, func(t trackPlayCount) model.Track { return t.track }), nil
}
