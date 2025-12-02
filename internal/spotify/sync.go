package spotify

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/pkg/storage"
	"github.com/topvennie/spotify_organizer/pkg/utils"
)

func (c *client) syncPlaylist(ctx context.Context, user model.User) (string, error) {
	playlistsDB, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return "", err
	}

	playlistsSpotify, err := c.playlistGetAll(ctx, user)
	if err != nil {
		return "", err
	}

	toCreate := make([]model.Playlist, 0)
	toUpdate := make([]model.Playlist, 0)
	toDelete := make([]model.Playlist, 0)

	// Find the playlists that need to be created or updated
	for i := range playlistsSpotify {
		playlistDB, ok := utils.SliceFind(playlistsDB, func(p *model.Playlist) bool { return p.Equal(playlistsSpotify[i]) })
		if !ok {
			// Playlist doesn't exist yet
			// Create it
			toCreate = append(toCreate, playlistsSpotify[i])
			continue
		}

		// Playlist already exist
		// But is it still completely the same?
		if !(*playlistDB).EqualEntry(playlistsSpotify[i]) {
			// Not completely the same anymore
			// Update it
			toUpdate = append(toUpdate, playlistsSpotify[i])
		}
	}

	// Do the database operations
	for i := range toCreate {
		if err := c.userCheck(ctx, toCreate[i].OwnerUID); err != nil {
			return "", err
		}
		if err := c.playlist.Create(ctx, &toCreate[i]); err != nil {
			return "", err
		}
	}

	for i := range toUpdate {
		if err := c.userCheck(ctx, toUpdate[i].OwnerUID); err != nil {
			return "", err
		}
		if err := c.playlist.Update(ctx, toUpdate[i]); err != nil {
			return "", err
		}
	}

	// New and updated entries are now in the database
	// Let's bring our local copy up to date
	playlistsDB, err = c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return "", err
	}

	// Find the playlists that need to be deleted
	for _, playlistDB := range playlistsDB {
		_, ok := utils.SliceFind(playlistsSpotify, func(p model.Playlist) bool { return p.SpotifyID == playlistDB.SpotifyID })
		if !ok {
			// Playlist no longer exists in the user's account
			// So delete it
			toDelete = append(toDelete, *playlistDB)
		}
	}

	for i := range toDelete {
		if err := c.playlist.Delete(ctx, toDelete[i].ID); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("Created: %d | Updated: %d | Deleted: %d", len(toCreate), len(toUpdate), len(toDelete)), nil
}

func (c *client) syncPlaylistCover(ctx context.Context, user model.User) (string, error) {
	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return "", err
	}
	if playlists == nil {
		return "", nil
	}

	newCovers := 0

	for _, playlist := range playlists {
		if playlist.CoverURL == "" {
			continue
		}

		cover, err := c.getCover(*playlist)
		if err != nil {
			return "", err
		}
		if len(cover) == 0 {
			continue
		}

		oldCover := []byte{}
		if playlist.CoverID != "" {
			oldCover, err = storage.S.Get(playlist.CoverID)
			if err != nil {
				return "", fmt.Errorf("get cover for %+v | %w", *playlist, err)
			}
		}

		if bytes.Equal(cover, oldCover) {
			continue
		}

		playlist.CoverID = uuid.NewString()
		if err := storage.S.Set(playlist.CoverID, cover, 0); err != nil {
			return "", fmt.Errorf("store new cover %+v | %w", *playlist, err)
		}

		if err := c.playlist.Update(ctx, *playlist); err != nil {
			return "", err
		}

		newCovers++
	}

	return fmt.Sprintf("New Covers: %d", newCovers), nil
}

// syncPlaylistTrack brings the local database up to date with the songs for each playlist
func (c *client) syncPlaylistTrack(ctx context.Context, user model.User) (string, error) {
	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return "", err
	}
	if playlists == nil {
		return "", nil
	}

	totalCreated := 0
	totalDeleted := 0

	for _, playlist := range playlists {
		tracksDB, err := c.track.GetByPlaylist(ctx, playlist.ID)
		if err != nil {
			return "", err
		}

		tracksSpotify, err := c.playlistGetTrackAll(ctx, user, *playlist)
		if err != nil {
			return "", err
		}

		toCreate := make([]model.Track, 0)
		toDelete := make([]model.Track, 0)

		for _, trackSpotify := range tracksSpotify {
			if _, ok := utils.SliceFind(tracksDB, func(t *model.Track) bool { return t.Equal(trackSpotify) }); !ok {
				toCreate = append(toCreate, trackSpotify)
			}
		}

		for _, trackDB := range tracksDB {
			if _, ok := utils.SliceFind(tracksSpotify, func(t model.Track) bool { return t.Equal(*trackDB) }); !ok {
				toDelete = append(toDelete, *trackDB)
			}
		}

		// Do the db operations
		for _, track := range toCreate {
			if err := c.trackCheck(ctx, &track); err != nil {
				return "", err
			}

			if err := c.playlist.CreateTrack(ctx, &model.PlaylistTrack{
				PlaylistID: playlist.ID,
				TrackID:    track.ID,
			}); err != nil {
				return "", err
			}
		}

		for _, track := range toDelete {
			if err := c.playlist.DeleteTrackByPlaylistTrack(ctx, model.PlaylistTrack{
				PlaylistID: playlist.ID,
				TrackID:    track.ID,
			}); err != nil {
				return "", err
			}
		}

		totalCreated += len(toCreate)
		totalDeleted += len(toDelete)
	}

	return fmt.Sprintf("Created %d | Deleted %d", totalCreated, totalDeleted), nil
}

func (c *client) syncLink(ctx context.Context, user model.User) (string, error) {
	directories, err := c.directory.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return "", err
	}

	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return "", err
	}

	links, err := c.link.GetAllByUser(ctx, user.ID)
	if err != nil {
		return "", err
	}

	totalAdded := 0

	for _, link := range links {
		var sources []model.Playlist
		var targets []model.Playlist

		switch {
		case link.SourceDirectoryID != 0:
			directory, ok := utils.SliceFind(directories, func(d *model.Directory) bool { return d.ID == link.SourceDirectoryID })
			if !ok {
				return "", fmt.Errorf("database foreign key reference error (source directory) for link %+v", *link)
			}
			sources = (*directory).Playlists

		case link.SourcePlaylistID != 0:
			playlist, ok := utils.SliceFind(playlists, func(p *model.Playlist) bool { return p.ID == link.SourcePlaylistID })
			if !ok {
				return "", fmt.Errorf("database foreign key reference error (source playlist) for link %+v", *link)
			}
			sources = []model.Playlist{**playlist}

		default:
			return "", fmt.Errorf("database foreign key reference error (source) for link %+v", *link)
		}

		switch {
		case link.TargetDirectoryID != 0:
			directory, ok := utils.SliceFind(directories, func(d *model.Directory) bool { return d.ID == link.TargetDirectoryID })
			if !ok {
				return "", fmt.Errorf("database foreign key reference error (target directory) for link %+v", *link)
			}
			targets = (*directory).Playlists

		case link.TargetPlaylistID != 0:
			playlist, ok := utils.SliceFind(playlists, func(p *model.Playlist) bool { return p.ID == link.TargetPlaylistID })
			if !ok {
				return "", fmt.Errorf("database foreign key reference error (target playlist) for link %+v", *link)
			}
			targets = []model.Playlist{**playlist}

		default:
			return "", fmt.Errorf("database foreign key reference error (target) for link %+v", *link)
		}

		for i := range sources {
			for j := range targets {
				added, err := c.syncLinkOne(ctx, user, sources[i], targets[j])
				if err != nil {
					return "", err
				}

				totalAdded += added
			}
		}
	}

	return fmt.Sprintf("Added %d", totalAdded), nil
}

func (c *client) syncLinkOne(ctx context.Context, user model.User, source, target model.Playlist) (int, error) {
	if source.Equal(target) {
		return 0, nil
	}

	tracksSource, err := c.track.GetByPlaylist(ctx, source.ID)
	if err != nil {
		return 0, err
	}

	tracksTarget, err := c.track.GetByPlaylist(ctx, target.ID)
	if err != nil {
		return 0o0, err
	}

	toAdd := make([]model.Track, 0)

	for _, trackSource := range tracksSource {
		if _, ok := utils.SliceFind(tracksTarget, func(t *model.Track) bool { return t.Equal(*trackSource) }); !ok {
			toAdd = append(toAdd, *trackSource)
		}
	}

	if err := c.playlistPostTrackAll(ctx, user, target, toAdd); err != nil {
		return 0, err
	}

	return len(toAdd), nil
}

// syncUser updates the information for every relevant user (for the given user)
func (c *client) syncUser(ctx context.Context, user model.User) error {
	// Get all relevant users
	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}
	if playlists == nil {
		return nil
	}

	usersDB := utils.SliceMap(playlists, func(p *model.Playlist) model.User { return p.Owner })
	usersDB = utils.SliceUnique(usersDB)

	// Get all spotify users
	usersSpotify := make([]model.User, 0, len(usersDB))
	for _, userDB := range usersDB {
		newUser, err := c.userGet(ctx, user, userDB)
		if err != nil {
			return err
		}

		usersSpotify = append(usersSpotify, newUser)
	}

	toUpdate := make([]model.User, 0)

	for _, userSpotify := range usersSpotify {
		if _, ok := utils.SliceFind(usersDB, func(u model.User) bool { return u.Equal(userSpotify) }); !ok {
			toUpdate = append(toUpdate, userSpotify)
		}
	}

	for _, user := range toUpdate {
		if err := c.user.Update(ctx, user); err != nil {
			return err
		}
	}

	return nil
}
