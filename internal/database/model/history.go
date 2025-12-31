package model

import (
	"time"

	"github.com/topvennie/sortifyr/pkg/sqlc"
)

type History struct {
	ID         int
	UserID     int
	TrackID    int
	PlayedAt   time.Time
	Skipped    *bool
	AlbumID    int
	ArtistID   int
	PlaylistID int
	ShowID     int

	// Non db fields
	Track     Track
	PlayCount int
}

func HistoryModel(h sqlc.History) *History {
	return &History{
		UserID:     int(h.UserID),
		ID:         int(h.ID),
		TrackID:    int(h.TrackID),
		PlayedAt:   h.PlayedAt.Time,
		Skipped:    fromBool(h.Skipped),
		AlbumID:    fromInt(h.AlbumID),
		ArtistID:   fromInt(h.ArtistID),
		PlaylistID: fromInt(h.PlaylistID),
		ShowID:     fromInt(h.ShowID),
	}
}

type HistoryFilter struct {
	UserID           int
	Limit            int
	Offset           int
	Start            time.Time
	End              time.Time
	Skipped          *bool
	PlayCountSkipped *bool
}
