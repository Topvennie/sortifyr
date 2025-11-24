package model

import (
	"time"

	"github.com/topvennie/spotify_organizer/pkg/sqlc"
)

type Playlist struct {
	ID            int
	SpotifyID     string
	OwnerID       int
	Name          string
	Description   string
	Public        bool
	Tracks        int
	Collaborative bool
	UpdatedAt     time.Time
	CreatedAt     time.Time

	// Non db fields
	Owner User
}

func PlaylistModelPopulated(p sqlc.Playlist, u sqlc.User) *Playlist {
	description := ""
	if p.Description.Valid {
		description = p.Description.String
	}

	return &Playlist{
		ID:            int(p.ID),
		SpotifyID:     p.SpotifyID,
		OwnerID:       int(p.OwnerID),
		Name:          p.Name,
		Description:   description,
		Public:        p.Public,
		Tracks:        int(p.Tracks),
		Collaborative: p.Collaborative,
		UpdatedAt:     p.UpdatedAt.Time,
		CreatedAt:     p.CreatedAt.Time,

		Owner: *UserModel(u),
	}
}

// Equal returns if 2 entries are equal (ignoring unique values)
func (p *Playlist) Equal(p2 Playlist) bool {
	return p.Owner.UID == p2.Owner.UID && p.Name == p2.Name && p.Description == p2.Description && p.Public == p2.Public && p.Tracks == p2.Tracks && p.Collaborative == p2.Collaborative
}
