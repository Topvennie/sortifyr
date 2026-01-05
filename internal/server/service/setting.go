package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
	"github.com/topvennie/sortifyr/internal/spotifysync"
	"github.com/topvennie/sortifyr/internal/task"
	"go.uber.org/zap"
)

const taskExportUID = "task-export"

var (
	exportNameReg  = regexp.MustCompile(`(?i)^.*_audio_.*\.json$`)
	exportIndexReg = regexp.MustCompile(`_(\d+).json$`)
)

type Setting struct {
	service Service

	history repository.History
	track   repository.Track
	user    repository.User
}

func (s *Service) NewSetting() *Setting {
	return &Setting{
		service: *s,
		history: *s.repo.NewHistory(),
		track:   *s.repo.NewTrack(),
		user:    *s.repo.NewUser(),
	}
}

func (s *Setting) Export(ctx context.Context, userID int, zip []byte) error {
	user, err := s.user.GetByID(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return fiber.ErrUnauthorized
	}

	if err := task.Manager.Add(ctx, task.NewTask(
		taskExportUID,
		"Import Spotify Export",
		task.IntervalOnce,
		false,
		func(ctx context.Context, _ []model.User) []task.TaskResult {
			return []task.TaskResult{{
				User:    *user,
				Message: "",
				Error:   s.exportTask(ctx, *user, zip),
			}}
		},
	)); err != nil {
		if errors.Is(err, task.ErrTaskExists) {
			return fiber.NewError(fiber.StatusBadRequest, "Task is already running")
		}
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

type exportTaskTrack struct {
	StoppedAt       time.Time `json:"ts"`
	Username        string    `json:"username"`
	MsPlayed        int       `json:"ms_played"`
	SpotifyTrackURI string    `json:"spotify_track_uri"`
	Skipped         bool      `json:"skipped"`
}

func (e exportTaskTrack) toHistory(userID int) *model.History {
	return &model.History{
		UserID:   userID,
		PlayedAt: e.StoppedAt.Add(time.Duration(-1*e.MsPlayed) * time.Millisecond),
		Skipped:  &e.Skipped,
		Track: model.Track{
			SpotifyID: spotifysync.UriToID(e.SpotifyTrackURI),
		},
	}
}

func (s *Setting) exportTask(ctx context.Context, user model.User, zipFile []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return fmt.Errorf("read zip file %w", err)
	}

	// Get all files
	files := make([]*zip.File, 0)

	for _, f := range zr.File {
		// Skip directories
		if f.FileInfo().IsDir() {
			continue
		}

		// Only read the audio files
		if match := exportNameReg.FindString(f.FileInfo().Name()); match == "" {
			continue
		}

		files = append(files, f)
	}
	if len(files) == 0 {
		return nil
	}

	// Go from the most recent to the oldest
	// Required for the exportFile function
	slices.SortFunc(files, func(a, b *zip.File) int {
		aIdx := exportIndex(a.FileInfo().Name())
		bIdx := exportIndex(b.FileInfo().Name())

		return bIdx - aIdx
	})

	for _, f := range files {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open file %s | %w", f.Name, err)
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("read file content %s | %w", f.FileInfo().Name(), err)
		}

		if err := s.exportTaskFile(ctx, user, content); err != nil {
			return err
		}
	}

	if err := task.Manager.RunRecurringByUID(spotifysync.TaskTrackUID, user); err != nil {
		return err
	}

	return nil
}

func (s *Setting) exportTaskFile(ctx context.Context, user model.User, file []byte) error {
	var exportTracks []exportTaskTrack
	if err := json.Unmarshal(file, &exportTracks); err != nil {
		return fmt.Errorf("parse file content to json %w", err)
	}

	// Get all entries
	exportHistory := make([]model.History, 0, len(exportTracks))
	spotifyIDs := make([]string, 0, len(exportTracks))
	for _, t := range exportTracks {
		if t.SpotifyTrackURI == "" {
			continue
		}

		h := *t.toHistory(user.ID)

		exportHistory = append(exportHistory, h)
		spotifyIDs = append(spotifyIDs, h.Track.SpotifyID)
	}
	slices.SortFunc(exportHistory, func(a, b model.History) int { return int(a.PlayedAt.UnixMilli() - b.PlayedAt.UnixMilli()) })

	if len(exportHistory) == 0 {
		return nil
	}

	// Get all already saved tracks
	tracksDB, err := s.track.GetAllBySpotify(ctx, spotifyIDs)
	if err != nil {
		return err
	}
	trackMap := make(map[string]int)
	for i := range tracksDB {
		trackMap[tracksDB[i].SpotifyID] = tracksDB[i].ID
	}

	// Populate db with any new track
	for i := range exportHistory {
		trackID, ok := trackMap[exportHistory[i].Track.SpotifyID]
		if !ok {
			// We don't have the track yet
			track := model.Track{SpotifyID: exportHistory[i].Track.SpotifyID}
			if err := s.track.Create(ctx, &track); err != nil {
				return err
			}
			trackID = track.ID
			trackMap[track.SpotifyID] = track.ID
		}

		exportHistory[i].TrackID = trackID
	}

	// Delete the old entries
	// This logic assumes that we go from the most recent track to the oldest
	if err := s.history.DeleteOlder(ctx, user.ID, exportHistory[len(exportHistory)-1].PlayedAt); err != nil {
		return err
	}
	if err := s.history.CreateBatch(ctx, exportHistory); err != nil {
		return err
	}

	return nil
}

func exportIndex(name string) int {
	matches := exportIndexReg.FindStringSubmatch(name)
	if len(matches) != 2 {
		return -1
	}

	idx, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1
	}

	return idx
}
