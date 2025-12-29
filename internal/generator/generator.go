// Package generator creates and maintaines playlists based on parameters and presets
package generator

import "github.com/topvennie/sortifyr/internal/database/model"

func Generate(preset model.GeneratorPreset, params model.GeneratorParameters) ([]model.Track, error) {
	var tracks []model.Track
	var err error

	switch preset {
	case model.PresetForgotten:
		tracks, err = forgotten(params)
	case model.PresetTop:
		tracks, err = top(params)
	case model.PresetOldTop:
		tracks, err = oldTop(params)

	default:
		tracks, err = custom(params)
	}

	if err != nil {
		return nil, err
	}

	return tracks, nil
}
