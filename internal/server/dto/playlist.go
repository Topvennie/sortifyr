package dto

import (
	"time"

	"github.com/topvennie/spotify_organizer/internal/database/model"
)

type Playlist struct {
	ID            int       `json:"id"`
	SpotifyID     string    `json:"spotify_id"`
	Owner         User      `json:"owner,omitzero"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitzero"`
	Public        bool      `json:"public"`
	Tracks        int       `json:"tracks"`
	Collaborative bool      `json:"collaborative"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func PlaylistDTO(playlist *model.Playlist, user *model.User) Playlist {
	return Playlist{
		ID:            playlist.ID,
		SpotifyID:     playlist.SpotifyID,
		Owner:         UserDTO(user),
		Name:          playlist.Name,
		Description:   playlist.Description,
		Public:        playlist.Public,
		Tracks:        playlist.Tracks,
		Collaborative: playlist.Collaborative,
		UpdatedAt:     playlist.UpdatedAt,
	}
}
