// Package generator creates and maintaines playlists based on parameters and presets
package generator

import (
	"context"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
)

type generator struct {
	history  repository.History
	playlist repository.Playlist
}

var G *generator

func Init(repo repository.Repository) {
	G = &generator{
		history:  *repo.NewHistory(),
		playlist: *repo.NewPlaylist(),
	}
}

func (g *generator) Create(ctx context.Context, gen model.Generator, playlist bool) error {
	// TODO: Implement
	return nil
}

func (g *generator) Edit(ctx context.Context, gen model.Generator, playlist bool) error {
	// TODO: IMplement
	return nil
}
