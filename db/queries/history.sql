-- name: HistoryGetPreviousPopulated :one
SELECT sqlc.embed(h), sqlc.embed(t)
FROM history h
LEFT JOIN tracks t ON t.id = h.track_id
WHERE h.played_at < $1 AND h.user_id = $2
ORDER BY h.played_at DESC
LIMIT 1;

-- name: HistoryGetPopulatedFiltered :many
SELECT sqlc.embed(h), sqlc.embed(t)
FROM history h
LEFT JOIN tracks t ON t.id = h.track_id
WHERE 
  h.user_id = $1::int AND
  (h.played_at >= $4::timestamptz OR NOT @filter_start) AND 
  (h.played_at <= $5::timestamptz OR NOT @filter_end) AND
  (h.skipped = $6 OR NOT @filter_skipped)
ORDER BY h.played_at DESC
LIMIT $2 OFFSET $3;
 
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

-- name: HistoryDeleteUserOlder :exec
DELETE FROM history
WHERE user_id = $1 AND played_at < $2;
