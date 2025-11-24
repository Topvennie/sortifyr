package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/internal/server/service"
)

type User struct {
	router fiber.Router

	user service.User
}

func NewUser(router fiber.Router, service service.Service) *User {
	api := &User{
		router: router.Group("/user"),
		user:   *service.NewUser(),
	}

	api.routes()

	return api
}

func (u *User) routes() {
	u.router.Get("/me", u.getMe)
}

func (u *User) getMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	user, err := u.user.GetByID(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(user)
}
