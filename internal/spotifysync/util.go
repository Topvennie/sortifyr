package spotifysync

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/topvennie/sortifyr/internal/spotifyapi"
	"github.com/topvennie/sortifyr/pkg/concurrent"
	"github.com/topvennie/sortifyr/pkg/storage"
)

func URIToID(uri string) string {
	parts := strings.Split(uri, ":")
	if len(parts) != 3 {
		return ""
	}

	return parts[2]
}

type syncUserDataStruct[T any] struct {
	DB             []T                 // All the items the user has in our db
	API            []T                 // All the items the user has according to spotify
	Equal          func(T, T) bool     // Test 2 items for equality
	Get            func(T) (*T, error) // Get a db item
	Create         func(*T) error      // Create a new db item
	CreateUserLink func(T) error       // Link a database item to the user
	DeleteUserLink func(T) error       // Delete the link between an item and the user
}

func syncUserData[T any](s syncUserDataStruct[T]) error {
	// Copy our database list
	// This allows us to delete items once found
	// which in turns supports multiple items in the same list
	// This is relevant for, for example, multiple instances of
	// the same track in a playlist
	dbCopy := make([]T, len(s.DB))
	copy(dbCopy, s.DB)

	// Go over every entry in the api list
	// to find missing entries
	for i := range s.API {
		if idx := slices.IndexFunc(dbCopy, func(t T) bool { return s.Equal(t, s.API[i]) }); idx != -1 {
			// The user already has this entry
			// Remove the entry from the list
			dbCopy[idx] = dbCopy[len(dbCopy)-1]
			dbCopy = dbCopy[:len(dbCopy)-1]

			continue
		}

		// The user doesn not have this item yet
		// Do we have it in our database?
		t, err := s.Get(s.API[i])
		if err != nil {
			return err
		}
		if t == nil {
			// We don't have it in our database yet
			// Let's create it so we get an id
			t = &s.API[i]
			if err := s.Create(t); err != nil {
				return err
			}

			// Add the item to the db list for the deletion loop
			s.DB = append(s.DB, *t)
		}

		// We now have an id for the db item
		// Let's link it to the user
		if err := s.CreateUserLink(*t); err != nil {
			return err
		}
	}

	// Same principle as the db copy but this time
	// to support deleting copies when the api no
	// longer has the copy
	apiCopy := make([]T, len(s.API))
	copy(apiCopy, s.API)

	// Do the same but let's look for
	// items that need to be deleted
	for i := range s.DB {
		if idx := slices.IndexFunc(apiCopy, func(t T) bool { return s.Equal(t, s.DB[i]) }); idx != -1 {
			// Item is in our api list
			// Let's delete it to support copies
			apiCopy[idx] = apiCopy[len(apiCopy)-1]
			apiCopy = apiCopy[:len(apiCopy)-1]

			continue
		}

		// Item is in our db but not in the spotify api list
		// So we have to delete it
		if err := s.DeleteUserLink(s.DB[i]); err != nil {
			return err
		}
	}

	return nil
}

type syncCoverStruct struct {
	CoverURL string
	CoverID  string
	Update   func(string) error
}

func (c *client) syncCover(ctx context.Context, s []syncCoverStruct) error {
	wg := concurrent.NewLimitedWaitGroup(12)

	var mu sync.Mutex
	var errs []error

	for _, item := range s {
		if item.CoverURL == "" {
			continue
		}

		wg.Go(func() {
			cover, err := spotifyapi.C.ImageGet(ctx, item.CoverURL)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			if len(cover) == 0 {
				return
			}

			oldCover := []byte{}
			if item.CoverID != "" {
				oldCover, err = storage.S.Get(item.CoverID)
				if err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("get cover for %+v | %w", item, err))
					mu.Unlock()
					return
				}
			}

			if bytes.Equal(cover, oldCover) {
				return
			}

			newID := uuid.NewString()
			if err := storage.S.Set(newID, cover, 0); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("store new cover %+v | %w", item, err))
				mu.Unlock()
				return
			}

			if err := item.Update(newID); err != nil {
				_ = storage.S.Delete(newID) // nolint:errcheck // Too bad if it fails
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			storage.S.Delete(item.CoverID) // nolint:errcheck // Too bad if it fails
		})
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

type filterSpotifyStruct[T any] struct {
	Items     []T                 // All items
	Frequency int                 // The frequency to update in hours
	SpotifyID func(t T) string    // Function to get the spotify id
	UpdatedAt func(t T) time.Time // Function to get the last time it was updated
}

func filterSpotify[T any](f filterSpotifyStruct[T]) []string {
	filtered := make([]string, 0)
	cutOff := time.Now().Add(time.Duration(-1*f.Frequency) * time.Hour)

	for i := range f.Items {
		if f.SpotifyID(f.Items[i]) == "" {
			continue
		}

		if f.UpdatedAt(f.Items[i]).After(cutOff) {
			continue
		}

		filtered = append(filtered, f.SpotifyID(f.Items[i]))
	}

	return filtered
}
