package model

import (
	"time"

	"github.com/topvennie/sortifyr/pkg/sqlc"
)

// JSON tags required for the generator repository

type Track struct {
	ID         int       `json:"id"`
	SpotifyID  string    `json:"spotify_id"`
	Name       string    `json:"name"`
	Popularity int       `json:"popularity"`
	DurationMs int       `json:"duration_ms"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Non db fields
	Artists   []Artist
	Playlist  Playlist
	CreatedAt time.Time
	DeletedAt time.Time
}

func TrackModel(t sqlc.Track) *Track {
	return &Track{
		ID:         int(t.ID),
		SpotifyID:  t.SpotifyID,
		Name:       fromString(t.Name),
		Popularity: fromInt(t.Popularity),
		DurationMs: fromInt(t.DurationMs),
		UpdatedAt:  fromTime(t.UpdatedAt),
	}
}

func (t *Track) Equal(t2 Track) bool {
	return t.SpotifyID == t2.SpotifyID
}

func (t *Track) EqualEntry(t2 Track) bool {
	return t.Name == t2.Name && t.Popularity == t2.Popularity && t.DurationMs == t2.DurationMs
}

type TrackArtist struct {
	ID       int
	TrackID  int
	ArtistID int
}

type TrackFilter struct {
	UserID     int
	PlaylistID int
	Limit      int
	Offset     int
}
