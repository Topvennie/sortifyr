package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/pkg/sqlc"
	"github.com/topvennie/spotify_organizer/pkg/utils"
)

type User struct {
	repo Repository
}

func (r *Repository) NewUser() *User {
	return &User{
		repo: *r,
	}
}

func (u *User) GetByID(ctx context.Context, id int) (*model.User, error) {
	user, err := u.repo.queries(ctx).UserGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user with id %d | %w", id, err)
	}

	return model.UserModel(user), nil
}

func (u *User) GetByUID(ctx context.Context, uid string) (*model.User, error) {
	user, err := u.repo.queries(ctx).UserGetByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user with uid %s | %w", uid, err)
	}

	return model.UserModel(user), nil
}

func (u *User) GetAllByID(ctx context.Context, ids []int) ([]*model.User, error) {
	users, err := u.repo.queries(ctx).UserGetAllByID(ctx, utils.SliceMap(ids, func(id int) int32 { return int32(id) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all users by id %v | %w", ids, err)
	}

	return utils.SliceMap(users, model.UserModel), nil
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	id, err := u.repo.queries(ctx).UserCreate(ctx, sqlc.UserCreateParams{
		Uid:         user.UID,
		Name:        user.Name,
		DisplayName: pgtype.Text{String: user.DisplayName, Valid: user.DisplayName != ""},
		Email:       user.Email,
	})
	if err != nil {
		return fmt.Errorf("create user %+v | %w", *user, err)
	}

	user.ID = int(id)

	return nil
}

func (u *User) Update(ctx context.Context, user model.User) error {
	if err := u.repo.queries(ctx).UserUpdate(ctx, sqlc.UserUpdateParams{
		ID:          int32(user.ID),
		Name:        user.Name,
		DisplayName: pgtype.Text{String: user.DisplayName, Valid: user.DisplayName != ""},
		Email:       user.Email,
	}); err != nil {
		return fmt.Errorf("update user %+v | %w", user, err)
	}

	return nil
}
