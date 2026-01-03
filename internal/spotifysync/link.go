package spotifysync

import (
	"context"
	"fmt"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/pkg/utils"
)

func (c *client) linksSync(ctx context.Context, user model.User) error {
	directories, err := c.directory.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	links, err := c.link.GetAllByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, link := range links {
		var sources []model.Playlist
		var targets []model.Playlist

		switch {
		case link.SourceDirectoryID != 0:
			directory, ok := utils.SliceFind(directories, func(d *model.Directory) bool { return d.ID == link.SourceDirectoryID })
			if !ok {
				return fmt.Errorf("database foreign key reference error (source directory) for link %+v", *link)
			}
			sources = (*directory).Playlists

		case link.SourcePlaylistID != 0:
			playlist, ok := utils.SliceFind(playlists, func(p *model.Playlist) bool { return p.ID == link.SourcePlaylistID })
			if !ok {
				return fmt.Errorf("database foreign key reference error (source playlist) for link %+v", *link)
			}
			sources = []model.Playlist{**playlist}

		default:
			return fmt.Errorf("database foreign key reference error (source) for link %+v", *link)
		}

		switch {
		case link.TargetDirectoryID != 0:
			directory, ok := utils.SliceFind(directories, func(d *model.Directory) bool { return d.ID == link.TargetDirectoryID })
			if !ok {
				return fmt.Errorf("database foreign key reference error (target directory) for link %+v", *link)
			}
			targets = (*directory).Playlists

		case link.TargetPlaylistID != 0:
			playlist, ok := utils.SliceFind(playlists, func(p *model.Playlist) bool { return p.ID == link.TargetPlaylistID })
			if !ok {
				return fmt.Errorf("database foreign key reference error (target playlist) for link %+v", *link)
			}
			targets = []model.Playlist{**playlist}

		default:
			return fmt.Errorf("database foreign key reference error (target) for link %+v", *link)
		}

		for i := range sources {
			for j := range targets {
				if err := c.linkOneSync(ctx, user, sources[i], targets[j]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *client) linkOneSync(ctx context.Context, user model.User, source, target model.Playlist) error {
	if source.Equal(target) {
		return nil
	}

	tracksSource, err := c.track.GetByPlaylist(ctx, source.ID)
	if err != nil {
		return err
	}

	tracksTarget, err := c.track.GetByPlaylist(ctx, target.ID)
	if err != nil {
		return err
	}

	toAdd := make([]model.Track, 0)

	for _, trackSource := range tracksSource {
		if _, ok := utils.SliceFind(tracksTarget, func(t *model.Track) bool { return t.Equal(*trackSource) }); !ok {
			toAdd = append(toAdd, *trackSource)
		}
	}

	if err := spotifyapi.C.PlaylistPostTrackAll(ctx, user, target.SpotifyID, toAdd); err != nil {
		return err
	}

	return nil
}
