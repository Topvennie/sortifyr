-- +goose Up
-- +goose StatementBegin
CREATE TABLE generators (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  playlist_id INTEGER REFERENCES playlists (id),
  maintained BOOL NOT NULL,
  interval INTEGER,
  parameters JSONB,

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
