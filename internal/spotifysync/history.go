package spotifysync

import (
	"context"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
)

func (c *client) historySync(ctx context.Context, user model.User) error {
	current, err := spotifyapi.C.PlayerGetCurrent(ctx, user)
	if err != nil {
		return err
	}
	if !current.IsPlaying {
		return nil
	}

	now := time.Now()
	currentStart := now.Add(time.Duration(-current.ProgressMs) * time.Millisecond)

	previous, err := c.history.GetPreviousPopulated(ctx, user.ID, now)
	if err != nil {
		return err
	}
	if previous == nil {
		previous = &model.History{}
	}

	if previous.Track.SpotifyID == current.Track.SpotifyID {
		// Same track
		// We don't know if the users is still listening for the first,
		// has it on repeat or simply paused it for a while.
		// Check first if it's still the same listening time (add a small buffer to count for e.g. network latency)
		if previous.PlayedAt.Add(10 * time.Second).After(currentStart) {
			// User is still in the same listen
			return nil
		}
		// Now the user has either paused it for some time or is listening to it on repeat.
		// Let's assume if the starting time of the current == the end time of the previous
		// that the user is listening on repeat.
		// To account for small time differences from crossplay, network latency, ... we'll againa add a buffer
		// but we need to add the buffer on both sides.
		previousEnd := previous.PlayedAt.Add(time.Duration(previous.Track.DurationMs) * time.Millisecond)
		if !previousEnd.Add(-20*time.Second).Before(currentStart) || !previousEnd.Add(20*time.Second).After(currentStart) {
			// User has paused for a while.
			// We could exit now and it would cover all cases except if
			// the user now starts listening to the song on repeat
			// To account for that we can move the playedAt time from the previous listen to the start
			// of the current track.
			// It's not perfect but it's the best we can do with the limited information
			previous.PlayedAt = currentStart
			return c.history.Update(ctx, *previous)
		}

		// We now know that the user is listening to the track on repeat, so let's add it again
	}

	track := current.Track.ToModel()
	if err := c.historyTrackCheck(ctx, &track); err != nil {
		return err
	}

	history := model.History{
		UserID:   user.ID,
		PlayedAt: currentStart,
		TrackID:  track.ID,
	}

	contextSpotifyID := uriToID(current.Context.URI)

	switch current.Context.Type {
	case "album":
		album := model.Album{SpotifyID: contextSpotifyID}
		if err := c.historyAlbumCheck(ctx, &album); err != nil {
			return err
		}
		history.AlbumID = album.ID
	case "artist":
		artist := model.Artist{SpotifyID: contextSpotifyID}
		if err := c.historyArtistCheck(ctx, &artist); err != nil {
			return err
		}
		history.ArtistID = artist.ID
	case "playlist":
		playlist := model.Playlist{SpotifyID: contextSpotifyID}
		if err := c.historyPlaylistCheck(ctx, &playlist); err != nil {
			return err
		}
		history.PlaylistID = playlist.ID
	case "show":
		show := model.Show{SpotifyID: contextSpotifyID}
		if err := c.historyShowCheck(ctx, &show); err != nil {
			return err
		}
		history.ShowID = show.ID
	}

	if err := c.history.Create(ctx, &history); err != nil {
		return err
	}

	return nil
}

func (c *client) historyTrackCheck(ctx context.Context, track *model.Track) error {
	trackDB, err := c.track.GetBySpotify(ctx, track.SpotifyID)
	if err != nil {
		return err
	}

	if trackDB == nil {
		if err := c.track.Create(ctx, track); err != nil {
			return err
		}
	} else {
		track.ID = trackDB.ID
	}

	return nil
}

func (c *client) historyArtistCheck(ctx context.Context, artist *model.Artist) error {
	artistDB, err := c.artist.GetBySpotify(ctx, artist.SpotifyID)
	if err != nil {
		return err
	}

	if artistDB == nil {
		if err := c.artist.Create(ctx, artist); err != nil {
			return err
		}
	} else {
		artist.ID = artistDB.ID
	}

	return nil
}

func (c *client) historyAlbumCheck(ctx context.Context, album *model.Album) error {
	albumDB, err := c.album.GetBySpotify(ctx, album.SpotifyID)
	if err != nil {
		return err
	}

	if albumDB == nil {
		if err := c.album.Create(ctx, album); err != nil {
			return err
		}
	} else {
		album.ID = albumDB.ID
	}

	return nil
}

func (c *client) historyPlaylistCheck(ctx context.Context, playlist *model.Playlist) error {
	playlistDB, err := c.playlist.GetBySpotify(ctx, playlist.SpotifyID)
	if err != nil {
		return err
	}

	if playlistDB == nil {
		if err := c.playlist.Create(ctx, playlist); err != nil {
			return err
		}
	} else {
		playlist.ID = playlistDB.ID
	}

	return nil
}

func (c *client) historyShowCheck(ctx context.Context, show *model.Show) error {
	showDB, err := c.show.GetBySpotify(ctx, show.SpotifyID)
	if err != nil {
		return err
	}

	if showDB == nil {
		if err := c.show.Create(ctx, show); err != nil {
			return err
		}
	} else {
		show.ID = showDB.ID
	}

	return nil
}

// historySkipped will populate the skipped field as much as it can
func (c *client) historySkipped(ctx context.Context, user model.User) error {
	histories, err := c.history.GetSkippedUnknownPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, h := range histories {
		previous, err := c.history.GetPreviousPopulated(ctx, user.ID, h.PlayedAt)
		if err != nil {
			return err
		}
		if previous == nil {
			continue
		}
		if previous.Track.DurationMs == 0 {
			continue
		}

		previousEnd := previous.PlayedAt.Add(time.Duration(previous.Track.DurationMs) * time.Millisecond)
		// Did the next song start before the previous ended?
		// Add a buffer of 20 seconds because skips in the last 20 seconds
		// of the track don't really count
		skipped := false
		if previousEnd.Add(-20 * time.Second).After(h.PlayedAt) {
			// User skipped
			skipped = true
		}

		previous.Skipped = &skipped
		if err := c.history.Update(ctx, *previous); err != nil {
			return err
		}
	}

	return nil
}
