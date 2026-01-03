-- name: GeneratorCreate :one
INSERT INTO generators (user_id, name, description, playlist_id, maintained, interval, parameters)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: GeneratorUpdate :exec
UPDATE generators
SET 
  name = coalesce(sqlc.narg('name'), name),
  description = coalesce(sqlc.narg('description'), description),
  playlist_id = coalesce(sqlc.narg('playlist_id'), playlist_id),
  maintained = coalesce(sqlc.narg('maintained'), maintained),
  interval = coalesce(sqlc.narg('interval'), interval),
  parameters = coalesce(sqlc.narg('parameters'), parameters)
WHERE id = $1;
