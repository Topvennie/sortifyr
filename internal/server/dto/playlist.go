package dto

import (
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type Playlist struct {
	ID            int    `json:"id" validate:"required"`
	SpotifyID     string `json:"spotify_id"`
	Owner         User   `json:"owner,omitzero"`
	Name          string `json:"name"`
	Description   string `json:"description,omitzero"`
	Public        *bool  `json:"public,omitzero"`
	TrackAmount   int    `json:"track_amount"`
	Collaborative *bool  `json:"collaborative,omitzero"`
	HasCover      bool   `json:"has_cover"`
}

func PlaylistDTO(playlist *model.Playlist, user *model.User) Playlist {
	return Playlist{
		ID:            playlist.ID,
		SpotifyID:     playlist.SpotifyID,
		Owner:         UserDTO(user),
		Name:          playlist.Name,
		Description:   playlist.Description,
		Public:        playlist.Public,
		TrackAmount:   playlist.TrackAmount,
		Collaborative: playlist.Collaborative,
		HasCover:      playlist.CoverID != "",
	}
}

func (p Playlist) ToModel() *model.Playlist {
	return &model.Playlist{
		ID:            p.ID,
		SpotifyID:     p.SpotifyID,
		OwnerID:       p.Owner.ID,
		Name:          p.Name,
		Description:   p.Description,
		Public:        p.Public,
		TrackAmount:   p.TrackAmount,
		Collaborative: p.Collaborative,
		Owner:         *p.Owner.ToModel(),
	}
}

type PlaylistDuplicate struct {
	Playlist
	Duplicates []TrackDuplicate `json:"duplicates"`
}

func PlaylistDuplicateDTO(playlist *model.Playlist, user *model.User, duplicates []model.Track) PlaylistDuplicate {
	type trackAmount struct {
		track  model.Track
		amount int
	}

	duplicateMap := make(map[int]trackAmount)
	for i := range duplicates {
		d, ok := duplicateMap[duplicates[i].ID]
		if !ok {
			d = trackAmount{track: duplicates[i], amount: 0}
		}

		d.amount++
		duplicateMap[duplicates[i].ID] = d
	}

	return PlaylistDuplicate{
		Playlist:   PlaylistDTO(playlist, user),
		Duplicates: utils.SliceMap(utils.MapValues(duplicateMap), func(t trackAmount) TrackDuplicate { return TrackDuplicateDTO(&t.track, t.amount) }),
	}
}

type PlaylistUnplayable struct {
	Playlist
	Unplayables []Track `json:"unplayables"`
}

func PlaylistUnplayableDTO(playlist *model.Playlist, user *model.User, unplayables []model.Track) PlaylistUnplayable {
	return PlaylistUnplayable{
		Playlist:    PlaylistDTO(playlist, user),
		Unplayables: utils.SliceMap(unplayables, func(t model.Track) Track { return TrackDTO(&t) }),
	}
}
