package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/internal/database/repository"
	"github.com/topvennie/spotify_organizer/internal/server/dto"
	"go.uber.org/zap"
)

type User struct {
	service Service

	user repository.User
}

func (s *Service) NewUser() *User {
	return &User{
		service: *s,
		user:    *s.repo.NewUser(),
	}
}

func (u *User) GetByID(ctx context.Context, id int) (dto.User, error) {
	user, err := u.user.GetByID(ctx, id)
	if err != nil {
		zap.S().Error(err)
		return dto.User{}, fiber.ErrInternalServerError
	}
	if user == nil {
		return dto.User{}, fiber.ErrNotFound
	}

	return dto.UserDTO(user), nil
}

func (u *User) GetByUID(ctx context.Context, uid string) (dto.User, error) {
	user, err := u.user.GetByUID(ctx, uid)
	if err != nil {
		zap.S().Error(err)
		return dto.User{}, fiber.ErrInternalServerError
	}
	if user == nil {
		return dto.User{}, fiber.ErrNotFound
	}

	return dto.UserDTO(user), nil
}

func (u *User) Create(ctx context.Context, userSave dto.User) (dto.User, error) {
	user := userSave.ToModel()

	if err := u.user.Create(ctx, user); err != nil {
		zap.S().Error(err)
		return dto.User{}, fiber.ErrInternalServerError
	}

	return dto.UserDTO(user), nil
}

func (u *User) Update(ctx context.Context, userSave dto.User) (dto.User, error) {
	user := userSave.ToModel()

	if err := u.user.Update(ctx, *user); err != nil {
		zap.S().Error(err)
		return dto.User{}, fiber.ErrInternalServerError
	}

	return dto.UserDTO(user), nil
}
