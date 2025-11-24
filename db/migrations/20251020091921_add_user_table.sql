
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  uid TEXT NOT NULL,
  name TEXT NOT NULL,
  display_name TEXT,
  email TEXT NOT NULL,

  UNIQUE (uid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

