-- name: TaskGetLastAllByUser :many
SELECT DISTINCT ON (task_uid) *
FROM task_runs
WHERE user_id = $1
ORDER BY task_uid, run_at DESC;

-- name: TaskRunCreate :one
INSERT INTO task_runs (task_uid, user_id, run_at, result, message, error, duration)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

