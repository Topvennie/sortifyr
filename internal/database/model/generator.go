package model

import (
	"encoding/json"
	"time"

	"github.com/topvennie/sortifyr/pkg/sqlc"
)

type GeneratorPreset string

const (
	GeneratorPresetCustom    GeneratorPreset = "custom"
	GeneratorPresetForgotten GeneratorPreset = "forgotten"
	GeneratorPresetTop       GeneratorPreset = "top"
	GeneratorPresetOldTop    GeneratorPreset = "old_top"
)

// We need json tags because the params are saved as jsonb

type GeneratorWindow struct {
	Start         time.Time     `json:"start"`
	End           time.Time     `json:"end"`
	MinPlays      int           `json:"min_plays"`
	BurstInterval time.Duration `json:"burst_interval"`
}

type GeneratorPresetCustomParams struct{}

type GeneratorPresetForgottenParams struct{}

type GeneratorPresetTopParams struct {
	Window GeneratorWindow `json:"window"`
}

type GeneratorPresetOldTopParams struct {
	PeakWindow   GeneratorWindow `json:"peak_window"`
	RecentWindow GeneratorWindow `json:"recent_window"`
}

type GeneratorParams struct {
	TrackAmount         int   `json:"track_amount"`
	ExcludedPlaylistIDs []int `json:"excluded_playlist_ids"`
	ExcludedTrackIDs    []int `json:"excluded_track_ids"`

	Preset GeneratorPreset `json:"preset"`

	ParamsCustom    *GeneratorPresetCustomParams    `json:"params_custom,omitzero"`
	ParamsForgotten *GeneratorPresetForgottenParams `json:"params_forgotten,omitzero"`
	ParamsTop       *GeneratorPresetTopParams       `json:"params_top,omitzero"`
	ParamsOldTop    *GeneratorPresetOldTopParams    `json:"params_old_top,omitzero"`
}

type Generator struct {
	ID              int
	UserID          int
	Name            string
	Description     string
	PlaylistID      int
	Interval        time.Duration
	SpotifyOutdated bool
	Params          GeneratorParams
	UpdatedAt       time.Time

	// Non db fields
	User   User
	Tracks []Track
}

func GeneratorModel(g sqlc.Generator) *Generator {
	params := GeneratorParams{}
	_ = json.Unmarshal(g.Parameters, &params) // nolint:errcheck // Data controlled by us, it'll be fine, ...right?

	return &Generator{
		ID:              int(g.ID),
		UserID:          int(g.UserID),
		Name:            g.Name,
		Description:     fromString(g.Description),
		PlaylistID:      fromInt(g.PlaylistID),
		Interval:        fromDuration(g.Interval),
		SpotifyOutdated: g.SpotifyOutdated,
		Params:          params,
		UpdatedAt:       g.UpdatedAt.Time,
	}
}

type GeneratorTrack struct {
	ID          int
	GeneratorID int
	TrackID     int
}
