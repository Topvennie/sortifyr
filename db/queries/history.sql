-- name: HistoryGetPreviousPopulated :one
SELECT sqlc.embed(h), sqlc.embed(t)
FROM history h
LEFT JOIN tracks t ON t.id = h.track_id
WHERE h.played_at < $1 AND h.user_id = $2
ORDER BY h.played_at DESC
LIMIT 1;

-- name: HistoryGetPopulatedFilteredPaginated :many
SELECT sqlc.embed(h), sqlc.embed(t), count(*) FILTER (WHERE h.user_id = $1::int AND (h.skipped = $7::boolean OR NOT @filter_play_count)) OVER  (PARTITION BY h.track_id) AS play_count
FROM history h
LEFT JOIN tracks t ON t.id = h.track_id
WHERE 
  h.user_id = $1::int AND
  (h.played_at >= $4::timestamptz OR NOT @filter_start) AND 
  (h.played_at <= $5::timestamptz OR NOT @filter_end) AND
  (h.skipped = $6::boolean OR NOT @filter_skipped)
ORDER BY h.played_at DESC
LIMIT $2 OFFSET $3;

-- name: HistoryGetPopulatedFiltered :many
SELECT sqlc.embed(h), sqlc.embed(t)
FROM history h
LEFT JOIN tracks t ON t.id = h.track_id
WHERE 
  h.user_id = $1::int AND
  (h.played_at >= $2::timestamptz OR NOT @filter_start) AND 
  (h.played_at <= $3::timestamptz OR NOT @filter_end) AND
  (h.skipped = $4::boolean OR NOT @filter_skipped)
ORDER BY h.played_at DESC;

-- name: HistoryGetSkippedNullPopulated :many
SELECT sqlc.embed(h), sqlc.embed(t)
FROM history h
LEFT JOIN tracks t ON t.id = h.track_id
WHERE h.skipped IS NULL AND h.user_id = $1
ORDER BY played_at ASC;
 
-- name: HistoryCreate :one
INSERT INTO history (user_id, track_id, played_at, album_id, artist_id, playlist_id, show_id, skipped)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: HistoryCreateBatch :exec
INSERT INTO history (user_id, track_id, played_at, skipped)
VALUES (
  UNNEST($1::int[]),
  UNNEST($2::int[]),
  UNNEST($3::timestamptz[]),
  UNNEST($4::boolean[])
);

-- name: HistoryUpdate :exec
UPDATE history
SET 
  played_at = coalesce(sqlc.narg('played_at'), played_at),
  skipped = coalesce(sqlc.narg('skipped'), skipped)
WHERE id = $1;

-- name: HistoryDeleteUserOlder :exec
DELETE FROM history
WHERE user_id = $1 AND played_at < $2;
