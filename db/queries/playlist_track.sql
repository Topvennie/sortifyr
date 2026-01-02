-- name: PlaylistTrackGetByPlaylistIds :many
SELECT *
FROM playlist_tracks
WHERE playlist_id = ANY($1::int[]);

-- name: PlaylistTrackCreate :one
INSERT INTO playlist_tracks (playlist_id, track_id)
VALUES ($1, $2)
RETURNING id;

-- name: PlaylistTrackDeleteByPlaylistTrack :exec
UPDATE playlist_tracks
SET deleted_at = NOW()
WHERE playlist_id = $1 AND track_id = $2;
