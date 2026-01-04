-- +goose Up
-- +goose StatementBegin
CREATE TABLE generators (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  playlist_id INTEGER REFERENCES playlists (id),
  maintained BOOLEAN NOT NULL,
  interval BIGINT,
  outdated BOOLEAN NOT NULL,
  parameters JSONB,
  updated_at TIMESTAMPTZ NOT NULL,

  CONSTRAINT generators_interval_required_if_maintained CHECK (
    maintained = false
    OR interval IS NOT NULL
  )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE generators;
-- +goose StatementEnd
