package spotifysync

import (
	"context"
	"errors"
	"fmt"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/task"
	"github.com/topvennie/sortifyr/pkg/config"
)

const (
	taskArtistUID            = "task-artist"
	taskAlbumUID             = "task-album"
	taskHistoryUID           = "task-history"
	taskLinkUID              = "task-link"
	taskPlaylistUID          = "task-playlist"
	taskPlaylistDuplicateUID = "task-playlist-duplicate"
	taskShowUID              = "task-show"
	taskTrackUID             = "task-track"
	taskUserUID              = "task-user"
	taskExportUID            = "task-export"
)

func (c *client) TaskPlaylistDuplicate(ctx context.Context, user model.User) error {
	if err := task.Manager.Add(ctx, task.NewTask(
		taskPlaylistDuplicateUID,
		"Playlist Duplicates Remove",
		task.IntervalOnce,
		func(ctx context.Context, _ []model.User) []task.TaskResult {
			results := []task.TaskResult{{
				User:    user,
				Message: "",
				Error:   nil,
			}}

			if err := c.playlistRemoveDuplicates(ctx, user); err != nil {
				results[0].Error = err
			}

			return results
		},
	)); err != nil {
		return err
	}

	return nil
}

func (c *client) TaskExport(ctx context.Context, user model.User, zip []byte) error {
	if err := task.Manager.Add(ctx, task.NewTask(
		taskExportUID,
		"Import Spotify Export",
		task.IntervalOnce,
		func(ctx context.Context, _ []model.User) []task.TaskResult {
			results := []task.TaskResult{{
				User:    user,
				Message: "",
				Error:   nil,
			}}

			if err := c.exportZip(ctx, user, zip); err != nil {
				results[0].Error = err
			}

			return results
		},
	)); err != nil {
		return err
	}

	return nil
}

func (c *client) taskRegister(ctx context.Context) error {
	if err := task.Manager.Add(ctx, task.NewTask(
		taskPlaylistUID,
		"Playlist",
		config.GetDefaultDuration("task.playlist_s", 60*60),
		c.taskWrap(c.taskPlaylist),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		taskAlbumUID,
		"Album",
		config.GetDefaultDuration("task.album_s", 60*60),
		c.taskWrap(c.taskAlbum),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		taskShowUID,
		"Show",
		config.GetDefaultDuration("task.show_s", 12*60*60),
		c.taskWrap(c.taskShow),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		taskTrackUID,
		"Track",
		config.GetDefaultDuration("task.track_s", 5*60),
		c.taskWrap(c.taskTrack),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		taskArtistUID,
		"Artist",
		config.GetDefaultDuration("task.artist_s", 5*60),
		c.taskWrap(c.taskArtist),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		taskUserUID,
		"User",
		config.GetDefaultDuration("task.user_s", 6*60*60),
		c.taskWrap(c.taskUser),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		taskHistoryUID,
		"Current",
		config.GetDefaultDuration("task.history_s", 15),
		c.taskWrap(c.taskHistory),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		taskLinkUID,
		"Link",
		config.GetDefaultDuration("task.link_s", 12*60*60),
		c.taskWrap(c.taskLink),
	)); err != nil {
		return err
	}

	return nil
}

func (c *client) taskWrap(fn func(context.Context, []model.User, []task.TaskResult)) func(context.Context, []model.User) []task.TaskResult {
	return func(ctx context.Context, users []model.User) []task.TaskResult {
		if len(users) == 0 {
			return []task.TaskResult{}
		}

		results := make([]task.TaskResult, 0, len(users))

		for _, user := range users {
			results = append(results, task.TaskResult{
				User:    user,
				Message: "",
				Error:   nil,
			})
		}

		fn(ctx, users, results)

		return results
	}
}

func (c *client) taskPlaylist(ctx context.Context, users []model.User, results []task.TaskResult) {
	for i, user := range users {
		if err := c.playlistSync(ctx, user); err != nil {
			results[i].Error = fmt.Errorf("synchronize playlists %w", err)
		}

		if err := c.playlistUpdate(ctx, user); err != nil {
			results[i].Error = errors.Join(fmt.Errorf("update playlists %w", err), results[i].Error)
		}

		if err := c.playlistCoverSync(ctx, user); err != nil {
			results[i].Error = errors.Join(fmt.Errorf("synchronize playlist covers %w", err), results[i].Error)
		}
	}
}

func (c *client) taskAlbum(ctx context.Context, users []model.User, results []task.TaskResult) {
	for i, user := range users {
		if err := c.albumSync(ctx, user); err != nil {
			results[i].Error = fmt.Errorf("synchronize albums %w", err)
		}

		if err := c.albumUpdate(ctx, user); err != nil {
			results[i].Error = errors.Join(fmt.Errorf("update albums %w", err), results[i].Error)
		}

		if err := c.albumCoverSync(ctx, user); err != nil {
			results[i].Error = errors.Join(fmt.Errorf("synchronize album covers %w", err), results[i].Error)
		}
	}
}

func (c *client) taskArtist(ctx context.Context, users []model.User, results []task.TaskResult) {
	if err := c.artistUpdate(ctx, users[0]); err != nil {
		for i := range users {
			results[i].Error = fmt.Errorf("update artists %w", err)
		}
	}
}

func (c *client) taskShow(ctx context.Context, users []model.User, results []task.TaskResult) {
	for i, user := range users {
		if err := c.showSync(ctx, user); err != nil {
			results[i].Error = fmt.Errorf("synchronize shows %w", err)
		}

		if err := c.showUpdate(ctx, user); err != nil {
			results[i].Error = errors.Join(fmt.Errorf("update shows %w", err), results[i].Error)
		}

		if err := c.showCoverSync(ctx, user); err != nil {
			results[i].Error = errors.Join(fmt.Errorf("synchronize shows covers %w", err), results[i].Error)
		}
	}
}

func (c *client) taskTrack(ctx context.Context, users []model.User, results []task.TaskResult) {
	if err := c.trackUpdate(ctx, users[0]); err != nil {
		for i := range users {
			results[i].Error = fmt.Errorf("update tracks %w", err)
		}
	}

	for i, user := range users {
		if err := c.historySkipped(ctx, user); err != nil {
			results[i].Error = errors.Join(fmt.Errorf("update historic data %w", err), results[i].Error)
		}
	}
}

func (c *client) taskUser(ctx context.Context, users []model.User, results []task.TaskResult) {
	for i, user := range users {
		if err := c.syncUser(ctx, user); err != nil {
			results[i].Error = fmt.Errorf("synchronize users %w", err)
		}
	}
}

func (c *client) taskHistory(ctx context.Context, users []model.User, results []task.TaskResult) {
	for i, user := range users {
		if err := c.historySync(ctx, user); err != nil {
			results[i].Error = fmt.Errorf("get history %w", err)
		}
	}
}

func (c *client) taskLink(ctx context.Context, users []model.User, results []task.TaskResult) {
	for i, user := range users {
		if err := c.linksSync(ctx, user); err != nil {
			results[i].Error = fmt.Errorf("synchronize links %w", err)
		}
	}
}
