-- name: GeneratorCreate :one
INSERT INTO generators (name, preset, parameters)
VALUES ($1, $2, $3)
RETURNING id;
