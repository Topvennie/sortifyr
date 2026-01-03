-- name: TrackGetAll :many
SELECT *
FROM tracks;

-- name: TrackGetAllById :many
SELECT *
FROM tracks
WHERE id = ANY($1::int[]);

-- name: TrackGetBySpotify :one
SELECT *
FROM tracks
WHERE spotify_id = $1;

-- name: TrackGetAllBySpotify :many
SELECT *
FROM tracks
WHERE spotify_id = ANY($1::text[]);

-- name: TrackGetByName :many
SELECT *
FROM tracks
WHERE name = $1;

-- name: TrackGetByPlaylist :many
SELECT t.*
FROM tracks t
LEFT JOIN playlist_tracks pt ON pt.track_id = t.id
WHERE pt.playlist_id = $1 AND pt.deleted_at IS NULL;

-- name: TrackGetCreatedFilteredPopulated :many
SELECT sqlc.embed(t), sqlc.embed(pt), sqlc.embed(p), sqlc.embed(u)
FROM tracks t
LEFT JOIN playlist_tracks pt ON pt.track_id = t.id
LEFT JOIN playlist_users pu ON pu.playlist_id = pt.playlist_id
LEFT JOIN playlists p ON p.id = pu.playlist_id
LEFT JOIN users u ON p.owner_id = u.id
WHERE 
  pu.user_id = $3::int AND 
  pt.deleted_at IS NULL AND
  (p.id = $4::int OR NOT @filter_playlist_id)
ORDER BY pt.created_at DESC
LIMIT $1 OFFSET $2;

-- name: TrackGetDeletedFilteredPopulated :many
SELECT sqlc.embed(t), sqlc.embed(pt), sqlc.embed(p), sqlc.embed(u)
FROM tracks t
LEFT JOIN playlist_tracks pt ON pt.track_id = t.id
LEFT JOIN playlist_users pu ON pu.playlist_id = pt.playlist_id
LEFT JOIN playlists p ON p.id = pu.playlist_id
LEFT JOIN users u ON p.owner_id = u.id
WHERE
  pu.user_id = $3::int AND 
  pt.deleted_at IS NOT NULL AND
  (p.id = $4::int OR NOT @filter_playlist_id)
ORDER BY pt.deleted_at DESC
LIMIT $1 OFFSET $2;

-- name: TrackCreate :one
INSERT INTO tracks (spotify_id, name, popularity, duration_ms)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: TrackUpdate :exec
UPDATE tracks
SET
  name = coalesce(sqlc.narg('name'), name),
  popularity = coalesce(sqlc.narg('popularity'), popularity),
  duration_ms = coalesce(sqlc.narg('duration_ms'), duration_ms),
  updated_at = NOW()
WHERE id = $1;
