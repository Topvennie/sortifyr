package model

import (
	"encoding/json"
	"time"

	"github.com/topvennie/sortifyr/pkg/sqlc"
)

type GeneratorPreset string

const (
	PresetCustom    GeneratorPreset = "custom"
	PresetForgotten GeneratorPreset = "forgotten"
	PresetTop       GeneratorPreset = "top"
	PresetOldTop    GeneratorPreset = "old_top"
)

type GeneratorWindow struct {
	Start         time.Time
	End           time.Time
	MinPlays      int
	BurstInterval time.Duration
}

type GeneratorParameters struct {
	UserID       int
	TrackAmount  int
	MinPlayCount int
	PeakWindow   GeneratorWindow
	RecentWindow GeneratorWindow
}

type Generator struct {
	ID         int
	Name       string
	Preset     GeneratorPreset
	Parameters GeneratorParameters
}

func GeneratorModel(g sqlc.Generator) *Generator {
	parameters := GeneratorParameters{}
	_ = json.Unmarshal(g.Parameters, &parameters)

	return &Generator{
		ID:         int(g.ID),
		Name:       g.Name,
		Preset:     GeneratorPreset(g.Preset),
		Parameters: parameters,
	}
}
