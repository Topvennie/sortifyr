package spotify

import (
	"context"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
)

func (c *client) historySync(ctx context.Context, user model.User) error {
	current, err := c.api.PlayerGetCurrent(ctx, user)
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
		// Let's give it a 5 second buffer
		if previous.PlayedAt.Add(5 * time.Second).After(currentStart) {
			return nil
		}
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
