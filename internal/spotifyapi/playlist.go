package spotifyapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/concurrent"
	"github.com/topvennie/sortifyr/pkg/utils"
)

func (c *client) PlaylistGet(ctx context.Context, user model.User, spotifyID string) (Playlist, error) {
	var resp Playlist

	if err := c.request(ctx, user, http.MethodGet, "playlists/"+spotifyID, http.NoBody, &resp); err != nil {
		return Playlist{}, fmt.Errorf("get playlist %s | %w", spotifyID, err)
	}

	return resp, nil
}

type playlistUserResponse struct {
	Total int        `json:"total"`
	Items []Playlist `json:"items"`
}

func (c *client) PlaylistGetUser(ctx context.Context, user model.User) ([]Playlist, error) {
	playlists := make([]Playlist, 0)

	total := 51
	limit := 50

	for i := 0; i < total; i += limit {
		var resp playlistUserResponse

		if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("me/playlists?offset=%d&limit=%d", i, limit), http.NoBody, &resp); err != nil {
			return nil, fmt.Errorf("get playlists with limit %d and offset %d | %w", limit, i, err)
		}

		playlists = append(playlists, resp.Items...)
		total = resp.Total
	}

	return playlists, nil
}

type playlistTrackAPI struct {
	Track Track `json:"track"`
}

type playlistTrackResponse struct {
	Total int                `json:"total"`
	Items []playlistTrackAPI `json:"items"`
}

func (c *client) PlaylistGetTrackAll(ctx context.Context, user model.User, spotifyID string) ([]Track, error) {
	wg := concurrent.NewLimitedWaitGroup(12)

	var mu sync.Mutex
	var errs []error

	tracks := make([]Track, 0)
	total := 0
	limit := 50

	// Do the first request to get the total
	var resp playlistTrackResponse
	if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("playlists/%s/tracks?offset=%d&limit=%d", spotifyID, 0, limit), http.NoBody, &resp); err != nil {
		return nil, fmt.Errorf("get playlist tracks with limit %d and offset %d | %w", limit, 0, err)
	}

	tracks = append(tracks, utils.SliceMap(resp.Items, func(t playlistTrackAPI) Track { return t.Track })...)
	total = resp.Total

	for i := limit; i < total; i += limit {
		wg.Go(func() {
			var resp playlistTrackResponse

			if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("playlists/%s/tracks?offset=%d&limit=%d", spotifyID, i, limit), http.NoBody, &resp); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("get playlist tracks with limit %d and offset %d | %w", limit, i, err))
				mu.Unlock()
				return
			}

			mu.Lock()
			tracks = append(tracks, utils.SliceMap(resp.Items, func(t playlistTrackAPI) Track { return t.Track })...)
			mu.Unlock()
		})
	}

	wg.Wait()

	return tracks, nil
}

type playlistTrackAddPayload struct {
	URIs []string `json:"uris"`
}

func (c *client) PlaylistPostTrackAll(ctx context.Context, user model.User, spotifyID string, tracks []model.Track) error {
	current := 0
	total := len(tracks)

	for current < total {
		end := min(current+100, total)

		payload := playlistTrackAddPayload{
			URIs: utils.SliceMap(tracks[current:end], func(t model.Track) string { return "spotify:track:" + t.SpotifyID }),
		}

		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal tracks add payload %+v | %w", payload, err)
		}

		body := bytes.NewReader(data)

		if err := c.request(ctx, user, http.MethodPost, fmt.Sprintf("playlists/%s/tracks", spotifyID), body, noResp); err != nil {
			return err
		}

		current = end
	}

	return nil
}

type playlistTrackRemovePayload struct {
	Tracks     []playlistTrackRemoveURIPayload `json:"tracks"`
	SnapshotID string                          `json:"snapshot_id"`
}

type playlistTrackRemoveURIPayload struct {
	URI string `json:"uri"`
}

func (c *client) PlaylistDeleteTrackAll(ctx context.Context, user model.User, spotifyID, snapshotID string, tracks []model.Track) error {
	current := 0
	total := len(tracks)

	for current < total {
		end := min(current+100, total)

		payload := playlistTrackRemovePayload{
			Tracks: utils.SliceMap(tracks[current:end], func(t model.Track) playlistTrackRemoveURIPayload {
				return playlistTrackRemoveURIPayload{URI: "spotify:track:" + t.SpotifyID}
			}),
			SnapshotID: snapshotID,
		}

		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal tracks delete payload %+v | %w", payload, err)
		}

		body := bytes.NewReader(data)

		if err := c.request(ctx, user, http.MethodDelete, fmt.Sprintf("playlists/%s/tracks", spotifyID), body, noResp); err != nil {
			return err
		}

		current = end
	}

	return nil
}

type playlistCreatePayload struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Public        bool   `json:"public"`
	Collaborative bool   `json:"collaborative"`
}

func (c *client) PlaylistCreate(ctx context.Context, user model.User, playlist *model.Playlist) error {
	payload := playlistCreatePayload{
		Name:          playlist.Name,
		Description:   playlist.Description,
		Public:        *playlist.Public,
		Collaborative: *playlist.Collaborative,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal create playlist payload %+v | %w", payload, err)
	}

	body := bytes.NewReader(data)

	var resp Playlist
	if err := c.request(ctx, user, http.MethodPost, fmt.Sprintf("/users/%s/playlists", user.UID), body, &resp); err != nil {
		return err
	}

	playlist.SpotifyID = resp.SpotifyID
	playlist.SnapshotID = resp.SnapshotID

	return nil
}

func (c *client) PlaylistDelete(ctx context.Context, user model.User, spotifyID string) error {
	return c.request(ctx, user, http.MethodDelete, fmt.Sprintf("/playlists/%s/followers", spotifyID), http.NoBody, noResp)
}
