-- +goose Up
-- +goose StatementBegin
ALTER TABLE task_runs
ADD COLUMN message TEXT;

DELETE FROM task_runs
WHERE user_id IS NULL; -- App is not in production yet so it's fine

ALTER TABLE task_runs
ALTER COLUMN user_id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE task_runs
ALTER COLUMN user_id DROP NOT NULL;

ALTER TABLE task_runs
DROP COLUMN message;
-- +goose StatementEnd
