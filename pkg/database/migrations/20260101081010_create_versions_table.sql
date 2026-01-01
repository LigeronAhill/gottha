-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS version(tag VARCHAR(31) PRIMARY KEY DEFAULT '0.0.1');

INSERT INTO
  version(tag)
VALUES
  ('v0.0.1');

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS version;

-- +goose StatementEnd
