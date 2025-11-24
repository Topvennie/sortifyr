-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlists (
  id SERIAL PRIMARY KEY,
  spotify_id TEXT NOT NULL,
  owner_id INTEGER NOT NULL REFERENCES users (id),
  name TEXT NOT NULL,
  description TEXT,
  public BOOLEAN NOT NULL,
  tracks INTEGER NOT NULL,
  collaborative BOOLEAN NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  UNIQUE (spotify_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE playlists;
-- +goose StatementEnd
