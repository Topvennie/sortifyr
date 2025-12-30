-- +goose Up
-- +goose StatementBegin
ALTER TABLE history
ADD COLUMN skipped BOOLEAN;

ALTER TABLE tracks
ADD COLUMN duration_ms INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tracks
DROP COLUMN duration_ms;

ALTER TABLE history
DROP COLUMN skipped;
-- +goose StatementEnd
