package dto

import "github.com/topvennie/sortifyr/internal/database/model"

type GeneratorParameters struct {
	TrackAmount  int `json:"track_amount" validate:"min=0"`
	MinPlayCount int `json:"min_play_count" validate:"min=0"`
}

func (g GeneratorParameters) ToModel() model.GeneratorParameters {
	return model.GeneratorParameters{
		TrackAmount:  g.TrackAmount,
		MinPlayCount: g.MinPlayCount,
	}
}

type Generator struct {
	Preset model.GeneratorPreset `json:"preset" validate:"required"`
	Params GeneratorParameters   `json:"params"`
}
