package spotifyapi

import (
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type Album struct {
	SpotifyID   string   `json:"id"`
	Name        string   `json:"name"`
	TrackAmount int      `json:"total_tracks"`
	Popularity  int      `json:"popularity"`
	Images      []Image  `json:"images"`
	Artists     []Artist `json:"artists"`
}

func (a Album) ToModel() model.Album {
	url := ""
	maxWidth := -1
	for _, image := range a.Images {
		if image.Width > maxWidth {
			url = image.URL
			maxWidth = image.Width
		}
	}

	return model.Album{
		SpotifyID:   a.SpotifyID,
		Name:        a.Name,
		TrackAmount: a.TrackAmount,
		Popularity:  a.Popularity,
		CoverURL:    url,
	}
}

type Artist struct {
	SpotifyID string `json:"id"`
	Name      string `json:"name"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Popularity int     `json:"popularity"`
	Images     []Image `json:"images"`
}

func (a Artist) ToModel() model.Artist {
	url := ""
	maxWidth := -1
	for _, image := range a.Images {
		if image.Width > maxWidth {
			url = image.URL
			maxWidth = image.Width
		}
	}

	return model.Artist{
		SpotifyID:  a.SpotifyID,
		Name:       a.Name,
		Followers:  a.Followers.Total,
		Popularity: a.Popularity,
		CoverURL:   url,
	}
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Context struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type History struct {
	Track    Track     `json:"track"`
	PlayedAt time.Time `json:"played_at"`
	Context  Context   `json:"context"`
}

type Current struct {
	Track      Track   `json:"item"`
	ProgressMs int     `json:"progress_ms"`
	IsPlaying  bool    `json:"is_playing"`
	Context    Context `json:"context"`
}

func (h History) ToModel(user model.User) model.History {
	return model.History{
		UserID:   user.ID,
		PlayedAt: h.PlayedAt,
	}
}

type Playlist struct {
	SpotifyID string `json:"id"`
	Owner     struct {
		UID         string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      *bool  `json:"public"`
	Tracks      struct {
		Total int `json:"total"`
	} `json:"tracks"`
	Collaborative *bool   `json:"collaborative"`
	Images        []Image `json:"images"`
	SnapshotID    string  `json:"snapshot_id"`
}

func (p *Playlist) ToModel() model.Playlist {
	url := ""
	maxWidth := -1
	for _, image := range p.Images {
		if image.Width > maxWidth {
			url = image.URL
			maxWidth = image.Width
		}
	}

	return model.Playlist{
		SpotifyID:     p.SpotifyID,
		Name:          p.Name,
		Description:   p.Description,
		Public:        p.Public,
		TrackAmount:   p.Tracks.Total,
		Collaborative: p.Collaborative,
		CoverURL:      url,
		Owner: model.User{
			UID:         p.Owner.UID,
			DisplayName: p.Owner.DisplayName,
		},
		SnapshotID: p.SnapshotID,
	}
}

type Show struct {
	SpotifyID     string  `json:"id"`
	Name          string  `json:"name"`
	EpisodeAmount int     `json:"total_episodes"`
	Images        []Image `json:"images"`
}

func (s Show) ToModel() model.Show {
	url := ""
	maxWidth := -1
	for _, image := range s.Images {
		if image.Width > maxWidth {
			url = image.URL
			maxWidth = image.Width
		}
	}

	return model.Show{
		SpotifyID:     s.SpotifyID,
		Name:          s.Name,
		EpisodeAmount: s.EpisodeAmount,
		CoverURL:      url,
	}
}

type Track struct {
	SpotifyID  string   `json:"id"`
	Name       string   `json:"name"`
	Popularity int      `json:"popularity"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	LinkedFrom struct {
		SpotifyID string `json:"id"`
	} `json:"linked_from"`
}

func (t *Track) ToModel() model.Track {
	spotifyID := t.SpotifyID
	if t.LinkedFrom.SpotifyID != "" {
		spotifyID = t.LinkedFrom.SpotifyID
	}

	return model.Track{
		SpotifyID:  spotifyID,
		Name:       t.Name,
		Popularity: t.Popularity,
		DurationMs: t.DurationMs,
		Artists:    utils.SliceMap(t.Artists, func(a Artist) model.Artist { return a.ToModel() }),
	}
}
