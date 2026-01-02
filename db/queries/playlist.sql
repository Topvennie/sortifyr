-- name: PlaylistGet :one
SELECT *
FROM playlists
WHERE id = $1;

-- name: PlaylistGetBySpotify :one
SELECT *
FROM playlists
WHERE spotify_id = $1;

-- name: PlaylistGetByUser :many
SELECT p.*
FROM playlists p
LEFT JOIN playlist_users pu ON pu.playlist_id = p.id
WHERE pu.user_id = $1 AND pu.deleted_at IS NULL
ORDER BY p.name;

-- name: PlaylistGetByUserWithOwner :many
SELECT sqlc.embed(p), sqlc.embed(u)
FROM playlists p
LEFT JOIN playlist_users pu ON pu.playlist_id = p.id
LEFT JOIN users u ON u.id = p.owner_id
WHERE pu.user_id = $1 AND p.owner_id IS NOT NULL AND pu.deleted_at IS NULL
ORDER BY p.name;

-- name: PlaylistGetDuplicateTracksByUser :many
SELECT sqlc.embed(p), sqlc.embed(t), sqlc.embed(u)
FROM playlist_tracks pt
JOIN (
  SELECT playlist_id, track_id
  FROM playlist_tracks
  WHERE deleted_at IS NULL
  GROUP BY playlist_id, track_id
  HAVING COUNT(*) > 1
) dup
ON dup.playlist_id = pt.playlist_id
AND dup.track_id = pt.track_id
LEFT JOIN playlists p ON p.id = pt.playlist_id
LEFT JOIN tracks t ON t.id = pt.track_id
LEFT JOIN playlist_users pu ON pu.playlist_id = p.id
LEFT JOIN users u ON u.id = p.owner_id
WHERE pu.user_id = $1 AND p.owner_id IS NOT NULL AND pu.deleted_at IS NULL
ORDER BY p.id, t.id, pt.id;

-- name: PlaylistGetUnplayableTracksByUser :many
SELECT sqlc.embed(p), sqlc.embed(t), sqlc.embed(u)
FROM playlist_tracks pt
LEFT JOIN playlists p ON p.id = pt.playlist_id
LEFT JOIN tracks t ON t.id = pt.track_id
LEFT JOIN playlist_users pu ON pu.playlist_id = p.id
LEFT JOIN users u ON u.id = p.owner_id
WHERE pu.user_id = $1 AND p.owner_id IS NOT NULL AND pu.deleted_at IS NULL AND (t.spotify_id IS NULL OR t.spotify_id = '')
ORDER BY p.id, t.id, pt.id;

-- name: PlaylistCreate :one
INSERT INTO playlists (spotify_id, owner_id, name, description, public, track_amount, collaborative, cover_id, cover_url, snapshot_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;

-- name: PlaylistUpdateBySpotify :exec
UPDATE playlists
SET 
  owner_id = coalesce(sqlc.narg('owner_id'), owner_id),
  name = coalesce(sqlc.narg('name'), name),
  description = coalesce(sqlc.narg('description'), description),
  public = coalesce(sqlc.narg('public'), public),
  track_amount = coalesce(sqlc.narg('track_amount'), track_amount),
  collaborative = coalesce(sqlc.narg('collaborative'), collaborative),
  cover_id = coalesce(sqlc.narg('cover_id'), cover_id),
  cover_url = coalesce(sqlc.narg('cover_url'), cover_url),
  snapshot_id = coalesce(sqlc.narg('snapshot_id'), snapshot_id),
  updated_at = NOW()
WHERE spotify_id = $1;

