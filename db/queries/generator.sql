-- name: GeneratorGet :one
SELECT *
FROM generators
WHERE id = $1;

-- name: GeneratorGetMaintainedPopulated :many
SELECT sqlc.embed(g), sqlc.embed(u)
FROM generators g
LEFT JOIN users u ON u.id = g.user_id
WHERE g.maintained = true;

-- name: GeneratorGetByUser :many
SELECT *
FROM generators
WHERE user_id = $1;

-- name: GeneratorCreate :one
INSERT INTO generators (user_id, name, description, playlist_id, maintained, interval, outdated, parameters, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
RETURNING id;

-- name: GeneratorUpdate :exec
UPDATE generators
SET 
  name = coalesce(sqlc.narg('name'), name),
  description = coalesce(sqlc.narg('description'), description),
  playlist_id = coalesce(sqlc.narg('playlist_id'), playlist_id),
  maintained = coalesce(sqlc.narg('maintained'), maintained),
  interval = coalesce(sqlc.narg('interval'), interval),
  outdated = coalesce(sqlc.narg('outdated'), outdated),
  parameters = coalesce(sqlc.narg('parameters'), parameters),
  updated_at = NOW()
WHERE id = $1;
