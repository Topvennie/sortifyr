package dto

import (
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
)

type Track struct {
	ID        int    `json:"id"`
	SpotifyID string `json:"spotify_id"`
	Name      string `json:"name"`
}

func TrackDTO(t *model.Track) Track {
	return Track{
		ID:        t.ID,
		SpotifyID: t.SpotifyID,
		Name:      t.Name,
	}
}

type TrackFilter struct {
	UserID     int
	PlaylistID int
	Limit      int
	Offset     int
}

func (t TrackFilter) ToModel() *model.TrackFilter {
	return &model.TrackFilter{
		UserID:     t.UserID,
		PlaylistID: t.PlaylistID,
		Limit:      t.Limit,
		Offset:     t.Offset,
	}
}

type TrackAdded struct {
	Track

	Playlist  Playlist  `json:"playlist"`
	CreatedAt time.Time `json:"created_at"`
}

func TrackAddedDTO(t *model.Track) TrackAdded {
	return TrackAdded{
		Track:     TrackDTO(t),
		Playlist:  PlaylistDTO(&t.Playlist, &t.Playlist.Owner),
		CreatedAt: t.CreatedAt,
	}
}

type TrackDeleted struct {
	Track

	Playlist  Playlist  `json:"playlist"`
	DeletedAt time.Time `json:"deleted_at"`
}

func TrackDeletedDTO(t *model.Track) TrackDeleted {
	return TrackDeleted{
		Track:     TrackDTO(t),
		Playlist:  PlaylistDTO(&t.Playlist, &t.Playlist.Owner),
		DeletedAt: t.DeletedAt,
	}
}

type TrackDuplicate struct {
	Track

	Amount int `json:"amount"`
}

func TrackDuplicateDTO(t *model.Track, amount int) TrackDuplicate {
	return TrackDuplicate{
		Track:  TrackDTO(t),
		Amount: amount,
	}
}

type History struct {
	Track

	HistoryID int       `json:"history_id"`
	PlayedAt  time.Time `json:"played_at"`
	PlayCount int       `json:"play_count,omitzero"`
}

func HistoryDTO(t *model.Track, h *model.History) History {
	return History{
		Track:     TrackDTO(t),
		HistoryID: h.ID,
		PlayedAt:  h.PlayedAt,
		PlayCount: h.PlayCount,
	}
}

type HistoryFilter struct {
	UserID  int
	Skipped *bool
	Start   time.Time
	End     time.Time
	Limit   int
	Offset  int
}

func (h HistoryFilter) ToModel() *model.HistoryFilter {
	return &model.HistoryFilter{
		UserID:  h.UserID,
		Skipped: h.Skipped,
		Start:   h.Start,
		End:     h.End,
		Limit:   h.Limit,
		Offset:  h.Offset,
	}
}
