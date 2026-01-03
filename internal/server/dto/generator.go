package dto

import (
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
)

type GeneratorWindow struct {
	Start          time.Time     `json:"start"`
	End            time.Time     `json:"end"`
	MinPlays       int           `json:"min_plays" validate:"min=0"`
	BurstIntervalS time.Duration `json:"burst_interval_s"`
}

func generatorWindowDTO(g model.GeneratorWindow) GeneratorWindow {
	return GeneratorWindow{
		Start:          g.Start,
		End:            g.End,
		MinPlays:       g.MinPlays,
		BurstIntervalS: g.BurstInterval / time.Second,
	}
}

func (g GeneratorWindow) ToModel() *model.GeneratorWindow {
	return &model.GeneratorWindow{
		Start:         g.Start,
		End:           g.End,
		MinPlays:      g.MinPlays,
		BurstInterval: g.BurstIntervalS * time.Second,
	}
}

type GeneratorPresetCustomParams struct{}

func generatorPresetCustomParamsDTO(params *model.GeneratorPresetCustomParams) *GeneratorPresetCustomParams {
	if params == nil {
		return nil
	}

	return &GeneratorPresetCustomParams{}
}

func (g GeneratorPresetCustomParams) ToModel() *model.GeneratorPresetCustomParams {
	return &model.GeneratorPresetCustomParams{}
}

type GeneratorPresetForgottenParams struct{}

func generatorPresetForgottenParamsDTO(params *model.GeneratorPresetForgottenParams) *GeneratorPresetForgottenParams {
	if params == nil {
		return nil
	}

	return &GeneratorPresetForgottenParams{}
}

func (g GeneratorPresetForgottenParams) ToModel() *model.GeneratorPresetForgottenParams {
	return &model.GeneratorPresetForgottenParams{}
}

type GeneratorPresetTopParams struct {
	Window GeneratorWindow `json:"window"`
}

func generatorPresetTopParamsDTO(params *model.GeneratorPresetTopParams) *GeneratorPresetTopParams {
	if params == nil {
		return nil
	}

	return &GeneratorPresetTopParams{
		Window: generatorWindowDTO(params.Window),
	}
}

func (g GeneratorPresetTopParams) ToModel() *model.GeneratorPresetTopParams {
	return &model.GeneratorPresetTopParams{
		Window: *g.Window.ToModel(),
	}
}

type GeneratorPresetOldTopParams struct {
	PeakWindow   GeneratorWindow `json:"peak_window"`
	RecentWindow GeneratorWindow `json:"recent_window"`
}

func generatorPresetOldTopParamsDTO(params *model.GeneratorPresetOldTopParams) *GeneratorPresetOldTopParams {
	if params == nil {
		return nil
	}

	return &GeneratorPresetOldTopParams{
		PeakWindow:   generatorWindowDTO(params.PeakWindow),
		RecentWindow: generatorWindowDTO(params.RecentWindow),
	}
}

func (g GeneratorPresetOldTopParams) ToModel() *model.GeneratorPresetOldTopParams {
	return &model.GeneratorPresetOldTopParams{
		PeakWindow:   *g.PeakWindow.ToModel(),
		RecentWindow: *g.RecentWindow.ToModel(),
	}
}

type GeneratorParams struct {
	TrackAmount         int   `json:"track_amount" validate:"min=0"`
	ExcludedPlaylistIDs []int `json:"excluded_playlist_ids,omitzero"`
	ExcludedTrackIDs    []int `json:"excluded_track_ids,omitzero"`

	Preset model.GeneratorPreset `json:"preset" validate:"required"`

	ParamsCustom    *GeneratorPresetCustomParams    `json:"params_custom,omitzero"`
	ParamsForgotten *GeneratorPresetForgottenParams `json:"params_forgotten,omitzero"`
	ParamsTop       *GeneratorPresetTopParams       `json:"params_top,omitzero"`
	ParamsOldTop    *GeneratorPresetOldTopParams    `json:"params_old_top,omitzero"`
}

func generatorParamsDTO(params model.GeneratorParams) GeneratorParams {
	return GeneratorParams{
		TrackAmount:         params.TrackAmount,
		ExcludedPlaylistIDs: params.ExcludedPlaylistIDs,
		ExcludedTrackIDs:    params.ExcludedTrackIDs,
		Preset:              params.Preset,
		ParamsCustom:        generatorPresetCustomParamsDTO(params.ParamsCustom),
		ParamsForgotten:     generatorPresetForgottenParamsDTO(params.ParamsForgotten),
		ParamsTop:           generatorPresetTopParamsDTO(params.ParamsTop),
		ParamsOldTop:        generatorPresetOldTopParamsDTO(params.ParamsOldTop),
	}
}

func (g GeneratorParams) ToModel() model.GeneratorParams {
	var paramsCustom *model.GeneratorPresetCustomParams
	if g.ParamsCustom != nil {
		paramsCustom = g.ParamsCustom.ToModel()
	}
	var paramsForgotten *model.GeneratorPresetForgottenParams
	if g.ParamsForgotten != nil {
		paramsForgotten = g.ParamsForgotten.ToModel()
	}
	var paramsTop *model.GeneratorPresetTopParams
	if g.ParamsTop != nil {
		paramsTop = g.ParamsTop.ToModel()
	}
	var paramsOldTop *model.GeneratorPresetOldTopParams
	if g.ParamsOldTop != nil {
		paramsOldTop = g.ParamsOldTop.ToModel()
	}
	return model.GeneratorParams{
		TrackAmount:         g.TrackAmount,
		ExcludedPlaylistIDs: g.ExcludedPlaylistIDs,
		ExcludedTrackIDs:    g.ExcludedTrackIDs,
		Preset:              g.Preset,
		ParamsCustom:        paramsCustom,
		ParamsForgotten:     paramsForgotten,
		ParamsTop:           paramsTop,
		ParamsOldTop:        paramsOldTop,
	}
}

type Generator struct {
	ID          int             `json:"id"`
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description,omitzero"`
	PlaylistID  int             `json:"playlist_id"`
	Maintained  bool            `json:"maintained"`
	IntervalS   time.Duration   `json:"interval_s"`
	Outdated    bool            `json:"outdated"`
	Params      GeneratorParams `json:"params" validate:"required"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func GeneratorDTO(gen *model.Generator) Generator {
	return Generator{
		ID:          gen.ID,
		Name:        gen.Name,
		Description: gen.Description,
		PlaylistID:  gen.PlaylistID,
		Maintained:  gen.Maintained,
		IntervalS:   gen.Interval / time.Second,
		Outdated:    gen.Outdated,
		Params:      generatorParamsDTO(gen.Params),
		UpdatedAt:   gen.UpdatedAt,
	}
}

type GeneratorSave struct {
	ID             int             `json:"id"`
	Name           string          `json:"name" validate:"required"`
	Description    string          `json:"description,omitzero"`
	CreatePlaylist bool            `json:"create_playlist"`
	Maintained     bool            `json:"maintained"`
	IntervalS      time.Duration   `json:"interval_s"`
	Params         GeneratorParams `json:"params" validate:"required"`
}

func (g GeneratorSave) ToModel(userID int) *model.Generator {
	return &model.Generator{
		ID:          g.ID,
		UserID:      userID,
		Name:        g.Name,
		Description: g.Description,
		Maintained:  g.Maintained,
		Interval:    g.IntervalS * time.Second,
		Params:      g.Params.ToModel(),
	}
}
