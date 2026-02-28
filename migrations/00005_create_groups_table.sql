-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS groups
(
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL CHECK (char_length(name) <= 255),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_groups_deleted_at ON groups(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_groups_deleted_at;
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd
