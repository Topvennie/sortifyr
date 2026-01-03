// Package spotifysync syncronizes our local database with the spotify api
package spotifysync

import (
	"context"

	"github.com/topvennie/sortifyr/internal/database/repository"
)

// Data updates are typically split in 3 parts
// Let's use albums as an example
//
// 	1. Get all the albums an user has saved.
// 		 Only save the album spotify id and link it to the user.
// 	2. Get all relevant albums for an user.
// 		 Get all the album data from spotify and sync the local db where needed
// 	3. Update the covers for each album.
//
// This makes more api calls then necessary but has clear seperations of responsibilities.
// It also simplifies the logic as spotify sometimes returns simplified objects when
// it is nested. For example when requesting the user's albums then you get
// simplified artist and track objects for each album. Or when you request all the user's playlists
// then you get an array of simplified playlist objects

type client struct {
	album     repository.Album
	artist    repository.Artist
	directory repository.Directory
	history   repository.History
	link      repository.Link
	playlist  repository.Playlist
	show      repository.Show
	track     repository.Track
	user      repository.User
}

var C *client

func Init(repo repository.Repository) error {
	C = &client{
		album:     *repo.NewAlbum(),
		artist:    *repo.NewArtist(),
		directory: *repo.NewDirectory(),
		history:   *repo.NewHistory(),
		link:      *repo.NewLink(),
		playlist:  *repo.NewPlaylist(),
		show:      *repo.NewShow(),
		track:     *repo.NewTrack(),
		user:      *repo.NewUser(),
	}

	if err := C.taskRegister(context.Background()); err != nil {
		return err
	}

	return nil
}
