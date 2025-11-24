package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/internal/server/service"
)

type Playlist struct {
	router fiber.Router

	playlist service.Playlist
}

func NewPlaylist(router fiber.Router, service service.Service) *Playlist {
	api := &Playlist{
		router:   router.Group("/playlist"),
		playlist: *service.NewPlaylist(),
	}

	api.routes()

	return api
}

func (p *Playlist) routes() {
	p.router.Get("/", p.getAll)
	p.router.Post("/sync", p.sync)
}

func (p *Playlist) getAll(c *fiber.Ctx) error {
	playlists, err := p.playlist.GetAll(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(playlists)
}

func (p *Playlist) sync(c *fiber.Ctx) error {
	uid, ok := c.Locals("spotifyID").(string)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := p.playlist.Sync(c.Context(), uid); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
