-- name: GeneratorGet :one
SELECT *
FROM generators
WHERE id = $1;

-- name: GeneratorGetAll :many
SELECT sqlc.embed(g), sqlc.embed(u)
FROM generators g
LEFT JOIN users u ON u.id = g.user_id;

-- name: GeneratorGetByUserPopulated :many
SELECT
  sqlc.embed(g),
  COALESCE(json_agg(t.*) FILTER (WHERE t.id IS NOT NULL), '[]')::jsonb AS tracks
FROM generators g
LEFT JOIN generator_tracks gt ON gt.generator_id = g.id
LEFT JOIN tracks t ON t.id = gt.track_id
WHERE g.user_id = $1
GROUP BY g.id;

-- name: GeneratorCreate :one
INSERT INTO generators (user_id, name, description, playlist_id, interval, spotify_outdated, parameters, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
RETURNING id;

-- name: GeneratorUpdate :exec
UPDATE generators
SET 
  name = coalesce(sqlc.narg('name'), name),
  description = coalesce(sqlc.narg('description'), description),
  playlist_id = sqlc.narg('playlist_id'),
  interval = coalesce(sqlc.narg('interval'), interval),
  spotify_outdated = coalesce(sqlc.narg('spotify_outdated'), spotify_outdated),
  parameters = coalesce(sqlc.narg('parameters'), parameters),
  updated_at = NOW()
WHERE id = $1;

-- name: GeneratorDelete :exec
DELETE FROM generators
WHERE id = $1;
