-- +goose Up
-- +goose StatementBegin
CREATE TYPE generator_preset AS ENUM ('custom', 'forgotten', 'top', 'old_top');

CREATE TABLE generators (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  preset GENERATOR_PRESET NOT NULL,
  parameters JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE generators;

DROP TYPE generator_preset;
-- +goose StatementEnd
