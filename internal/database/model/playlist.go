package model

import (
	"time"

	"github.com/topvennie/sortifyr/pkg/sqlc"
)

type Playlist struct {
	ID            int
	SpotifyID     string
	OwnerID       int
	Name          string
	Description   string
	Public        *bool
	TrackAmount   int
	Collaborative *bool
	CoverID       string
	CoverURL      string
	SnapshotID    string
	UpdatedAt     time.Time

	// Non db fields
	Owner       User
	Duplicates  []Track
	Unplayables []Track
}

func PlaylistModel(p sqlc.Playlist) *Playlist {
	return &Playlist{
		ID:            int(p.ID),
		SpotifyID:     p.SpotifyID,
		OwnerID:       fromInt(p.OwnerID),
		Name:          fromString(p.Name),
		Description:   fromString(p.Description),
		Public:        fromBool(p.Public),
		TrackAmount:   fromInt(p.TrackAmount),
		Collaborative: fromBool(p.Collaborative),
		CoverID:       fromString(p.CoverID),
		CoverURL:      fromString(p.CoverUrl),
		SnapshotID:    fromString(p.SnapshotID),
		UpdatedAt:     fromTime(p.UpdatedAt),
	}
}

func PlaylistModelPopulated(p sqlc.Playlist, u sqlc.User) *Playlist {
	playlist := PlaylistModel(p)
	playlist.Owner = *UserModel(u)

	return playlist
}

func (p *Playlist) Equal(p2 Playlist) bool {
	return p.SpotifyID == p2.SpotifyID
}

func (p *Playlist) EqualEntry(p2 Playlist) bool {
	return p.OwnerID == p2.OwnerID && p.Name == p2.Name && p.Description == p2.Description && p.Public == p2.Public && p.TrackAmount == p2.TrackAmount && p.Collaborative == p2.Collaborative && p.CoverURL == p2.CoverURL && p.SnapshotID == p2.SnapshotID
}

type PlaylistTrack struct {
	ID         int
	PlaylistID int
	TrackID    int
	CreatedAt  time.Time
	DeletedAt  time.Time
}

type PlaylistUser struct {
	ID         int
	UserID     int
	PlaylistID int
	DeletedAt  time.Time
}
