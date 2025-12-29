package model

import (
	"encoding/json"

	"github.com/topvennie/sortifyr/pkg/sqlc"
)

type GeneratorPreset string

const (
	PresetCustom    GeneratorPreset = "custom"
	PresetForgotten GeneratorPreset = "forgotten"
	PresetTop       GeneratorPreset = "top"
	PresetOldTop    GeneratorPreset = "old_top"
)

type GeneratorParameters struct{}

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
