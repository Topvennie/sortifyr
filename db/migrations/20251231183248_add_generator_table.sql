-- +goose Up
-- +goose StatementBegin
CREATE TABLE generators (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  playlist_id INTEGER REFERENCES playlists (id),
  interval BIGINT,
  spotify_outdated BOOLEAN NOT NULL,
  parameters JSONB,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE generator_tracks (
  id SERIAL PRIMARY KEY,
  generator_id INTEGER NOT NULL REFERENCES generators (id) ON DELETE CASCADE,
  track_id INTEGER NOT NULL REFERENCES tracks (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE generator_tracks;

DROP TABLE generators;
-- +goose StatementEnd
