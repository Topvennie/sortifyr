-- name: PlaylistGetAllWithOwner :many
SELECT sqlc.embed(p), sqlc.embed(u)
FROM playlists p
LEFT JOIN users u ON u.id = p.owner_id
ORDER BY p.name;

-- name: PlaylistCreate :one
INSERT INTO playlists (spotify_id, owner_id, name, description, public, tracks, collaborative)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: PlaylistUpdate :exec
UPDATE playlists
SET owner_id = $2, name = $3, description = $4, public = $5, tracks = $6, collaborative = $7, updated_at = NOW()
WHERE id = $1;

-- name: PlaylistDelete :exec
DELETE FROM playlists
WHERE id = $1;
