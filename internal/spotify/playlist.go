// Package spotify connects with the spotify API
package spotify

import (
	"context"
	"fmt"

	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/pkg/utils"
)

func (c *client) PlaylistSync(ctx context.Context, uid string) error {
	user, err := c.user.GetByUID(ctx, uid)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found %s", uid)
	}

	playlistsDB, err := c.playlist.GetAllPopulated(ctx)
	if err != nil {
		return err
	}

	playlistsSpotify, err := c.GetAllPlaylists(ctx, uid)
	if err != nil {
		return err
	}

	toCreate := make([]model.Playlist, 0)
	toUpdate := make([]model.Playlist, 0)
	toDelete := make([]model.Playlist, 0)

	// Find the playlists that need to be created or updated
	for i := range playlistsSpotify {
		playlistDB, ok := utils.SliceFind(playlistsDB, func(p *model.Playlist) bool { return p.SpotifyID == playlistsSpotify[i].SpotifyID })
		if !ok {
			// Playlist doesn't exist yet
			// Create it
			toCreate = append(toCreate, playlistsSpotify[i])
			continue
		}

		// Playlist already exist
		// But is it still completely the same?
		if !(*playlistDB).Equal(playlistsSpotify[i]) {
			// Not completely the same anymore
			// Update it
			toUpdate = append(toUpdate, playlistsSpotify[i])
		}
	}

	for i := range toCreate {
		playlistOwner, err := c.user.GetByUID(ctx, toCreate[i].Owner.UID)
		if err != nil {
			return err
		}
		if playlistOwner == nil {
			if err := c.user.Create(ctx, &toCreate[i].Owner); err != nil {
				return err
			}
			playlistOwner = &toCreate[i].Owner
		} else if playlistOwner.DisplayName != toCreate[i].Owner.DisplayName {
			toCreate[i].Owner.ID = playlistOwner.ID
			if err := c.user.Update(ctx, toCreate[i].Owner); err != nil {
				return err
			}
			playlistOwner = &toCreate[i].Owner
		}
		toCreate[i].OwnerID = playlistOwner.ID

		if err := c.playlist.Create(ctx, &toCreate[i]); err != nil {
			return err
		}
	}

	for i := range toUpdate {
		playlistOwner, err := c.user.GetByUID(ctx, toUpdate[i].Owner.UID)
		if err != nil {
			return err
		}
		if playlistOwner == nil {
			if err := c.user.Create(ctx, &toUpdate[i].Owner); err != nil {
				return err
			}
			playlistOwner = &toUpdate[i].Owner
		} else if playlistOwner.DisplayName != toUpdate[i].Owner.DisplayName {
			toUpdate[i].Owner.ID = playlistOwner.ID
			if err := c.user.Update(ctx, toUpdate[i].Owner); err != nil {
				return err
			}
			playlistOwner = &toUpdate[i].Owner
		}
		toUpdate[i].OwnerID = playlistOwner.ID

		if err := c.playlist.Update(ctx, toUpdate[i]); err != nil {
			return err
		}
	}

	// New and updated entries are now in the database
	// Let's bring our local copy up to date
	playlistsDB, err = c.playlist.GetAllPopulated(ctx)
	if err != nil {
		return err
	}

	// Find the playlists that need to be deleted
	for _, playlistDB := range playlistsDB {
		_, ok := utils.SliceFind(playlistsSpotify, func(p model.Playlist) bool { return p.SpotifyID == playlistDB.SpotifyID })
		if !ok {
			// Playlist no longer exists in the user's account
			// So delete it
			toDelete = append(toDelete, *playlistDB)
		}
	}

	for i := range toDelete {
		if err := c.playlist.Delete(ctx, toDelete[i].ID); err != nil {
			return err
		}
	}

	return nil
}

type playlist struct {
	SpotifyID string `json:"id"`
	Owner     struct {
		UID         string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Tracks      struct {
		Total int `json:"total"`
	} `json:"tracks"`
	Collaborative bool `json:"collaborative"`
}

func (p *playlist) toModel() *model.Playlist {
	return &model.Playlist{
		SpotifyID:     p.SpotifyID,
		Name:          p.Name,
		Description:   p.Description,
		Public:        p.Public,
		Tracks:        p.Tracks.Total,
		Collaborative: p.Collaborative,
		Owner: model.User{
			UID:         p.Owner.UID,
			DisplayName: p.Owner.DisplayName,
		},
	}
}

type playListResponse struct {
	Total int        `json:"total"`
	Items []playlist `json:"items"`
}

func (c *client) GetAllPlaylists(ctx context.Context, uid string) ([]model.Playlist, error) {
	playlists := make([]model.Playlist, 0)

	limit := 50
	offset := 0

	resp, err := c.getPlaylists(ctx, uid, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get playlist with limit %d and offset %d | %w", limit, offset, err)
	}
	playlists = append(playlists, utils.SliceMap(resp.Items, func(p playlist) model.Playlist { return *p.toModel() })...)

	total := resp.Total

	for offset+limit < total {
		offset += limit

		resp, err := c.getPlaylists(ctx, uid, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("get playlist with limit %d and offset %d | %w", limit, offset, err)
		}
		playlists = append(playlists, utils.SliceMap(resp.Items, func(p playlist) model.Playlist { return *p.toModel() })...)
	}

	return playlists, nil
}

func (c *client) getPlaylists(ctx context.Context, uid string, limit, offset int) (playListResponse, error) {
	var resp playListResponse

	if err := c.request(ctx, uid, fmt.Sprintf("me/playlists?offset=%d&limit=%d", offset, limit), &resp); err != nil {
		return resp, fmt.Errorf("get playlist %w", err)
	}

	return resp, nil
}
