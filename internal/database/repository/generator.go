package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/sqlc"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type Generator struct {
	repo Repository
}

func (r *Repository) NewGenerator() *Generator {
	return &Generator{
		repo: *r,
	}
}

func (g *Generator) Get(ctx context.Context, id int) (*model.Generator, error) {
	gen, err := g.repo.queries(ctx).GeneratorGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get generator by id %d | %w", id, err)
	}

	return model.GeneratorModel(gen), nil
}

func (g *Generator) GetByUser(ctx context.Context, userID int) ([]*model.Generator, error) {
	gens, err := g.repo.queries(ctx).GeneratorGetByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get generators by user %d | %w", userID, err)
	}

	return utils.SliceMap(gens, model.GeneratorModel), nil
}

func (g *Generator) Create(ctx context.Context, gen *model.Generator) error {
	params, err := json.Marshal(gen.Params)
	if err != nil {
		return fmt.Errorf("create generator marshal params %+v | %w", *gen, err)
	}

	id, err := g.repo.queries(ctx).GeneratorCreate(ctx, sqlc.GeneratorCreateParams{
		UserID:      int32(gen.UserID),
		Name:        gen.Name,
		Description: toString(gen.Description),
		PlaylistID:  toInt(gen.PlaylistID),
		Maintained:  gen.Maintained,
		Interval:    toDuration(gen.Interval),
		Outdated:    gen.Outdated,
		Parameters:  params,
	})
	if err != nil {
		return fmt.Errorf("create generator %+v | %w", *gen, err)
	}

	gen.ID = int(id)

	return nil
}

func (g *Generator) Update(ctx context.Context, gen model.Generator) error {
	params, err := json.Marshal(gen.Params)
	if err != nil {
		return fmt.Errorf("update generator marshal params %+v | %w", gen, err)
	}

	if err := g.repo.queries(ctx).GeneratorUpdate(ctx, sqlc.GeneratorUpdateParams{
		ID:          int32(gen.ID),
		Name:        toString(gen.Name),
		Description: toString(gen.Description),
		PlaylistID:  toInt(gen.PlaylistID),
		Maintained:  toBool(&gen.Maintained),
		Interval:    toDuration(gen.Interval),
		Outdated:    toBool(&gen.Outdated),
		Parameters:  params,
	}); err != nil {
		return fmt.Errorf("update generator %+v | %w", gen, err)
	}

	return nil
}
