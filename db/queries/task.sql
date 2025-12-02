-- name: TaskGetByUID :one
SELECT *
FROM tasks
WHERE uid = $1;

-- name: TaskRunGet :one
SELECT sqlc.embed(t), sqlc.embed(r)
FROM task_runs r
LEFT JOIN tasks t ON t.uid = r.task_uid
WHERE r.id = $1;

-- name: TaskGetFiltered :many
SELECT sqlc.embed(t), sqlc.embed(r)
FROM task_runs r
LEFT JOIN tasks t ON t.uid = r.task_uid
WHERE
  r.user_id = $1::int AND
  (t.uid = $2 OR NOT @filter_task_uid) AND
  (r.result = $3 OR NOT @filter_result) AND
  t.active
ORDER BY r.run_at DESC
LIMIT $4 OFFSET $5;

-- name: TaskCreate :exec
INSERT INTO tasks (uid, name, active)
VALUES ($1, $2, $3);

-- name: TaskUpdate :exec
UPDATE tasks
SET name = $2, active = $3
WHERE uid = $1;

-- name: TaskSetInactiveAll :exec
UPDATE tasks
SET active = false;
