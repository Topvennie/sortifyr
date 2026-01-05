package generator

import (
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
)

// normalize sets default values for the parameters
func normalize(gen *model.Generator) {
	params := gen.Params
	if params.TrackAmount == 0 {
		params.TrackAmount = 50
	}

	switch params.Preset {
	case model.GeneratorPresetCustom:
		normalizePresetCustom(&params)
		params.ParamsForgotten = nil
		params.ParamsTop = nil
		params.ParamsOldTop = nil
	case model.GeneratorPresetForgotten:
		normalizePresetForgotten(&params)
		params.ParamsCustom = nil
		params.ParamsTop = nil
		params.ParamsOldTop = nil
	case model.GeneratorPresetTop:
		normalizePresetTop(&params)
		params.ParamsCustom = nil
		params.ParamsForgotten = nil
		params.ParamsOldTop = nil
	case model.GeneratorPresetOldTop:
		normalizePresetOldTop(&params)
		params.ParamsCustom = nil
		params.ParamsForgotten = nil
		params.ParamsTop = nil
	}

	gen.Params = params
}

func normalizePresetCustom(params *model.GeneratorParams) {
}

func normalizePresetForgotten(params *model.GeneratorParams) {
}

func normalizePresetTop(params *model.GeneratorParams) {
	now := time.Now()

	defaultParams := model.GeneratorPresetTopParams{
		Window: model.GeneratorWindow{
			Start:         now.Add(-14 * 24 * time.Hour), // 14 days ago
			End:           now,
			MinPlays:      5,
			BurstInterval: 14 * 24 * time.Hour, // 14 days
		},
	}

	if params.ParamsTop == nil {
		params.ParamsTop = &defaultParams
		return
	}

	normalizeWindow(&params.ParamsTop.Window, defaultParams.Window)
}

func normalizePresetOldTop(params *model.GeneratorParams) {
	now := time.Now()

	defaultParams := model.GeneratorPresetOldTopParams{
		PeakWindow: model.GeneratorWindow{
			Start:         now.Add(-1 * 24 * 365 * time.Hour), // 365 days ago
			End:           now.Add(-1 * 24 * 100 * time.Hour), // 100 days ago
			MinPlays:      5,
			BurstInterval: 24 * 14 * time.Hour, // 14 days
		},
		RecentWindow: model.GeneratorWindow{
			Start:         now.Add(-1 * 24 * 30 * time.Hour), // 30 days ago
			End:           now,                               // now
			MinPlays:      2,
			BurstInterval: 24 * 14 * time.Hour, // 14 days
		},
	}

	if params.ParamsOldTop == nil {
		params.ParamsOldTop = &defaultParams
		return
	}

	normalizeWindow(&params.ParamsOldTop.PeakWindow, defaultParams.PeakWindow)
	normalizeWindow(&params.ParamsOldTop.RecentWindow, defaultParams.RecentWindow)
}

func normalizeWindow(window *model.GeneratorWindow, normalized model.GeneratorWindow) {
	if window.Start.IsZero() {
		window.Start = normalized.Start
	}
	if window.End.IsZero() {
		window.End = normalized.End
	}
	if window.End.Before(window.Start) {
		window.End = window.Start
	}

	if window.MinPlays == 0 {
		window.MinPlays = normalized.MinPlays
	}

	if window.BurstInterval == 0 {
		window.BurstInterval = normalized.BurstInterval
	}
}
