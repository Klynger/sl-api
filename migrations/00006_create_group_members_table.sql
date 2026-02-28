-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS group_members
(
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  group_id UUID NOT NULL,
  roles TEXT[] DEFAULT ARRAY['member']::TEXT[] NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
  CONSTRAINT unique_user_group UNIQUE (user_id, group_id)
);
CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members(user_id);
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_deleted_at ON group_members(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_group_members_deleted_at;
DROP INDEX IF EXISTS idx_group_members_group_id;
DROP INDEX IF EXISTS idx_group_members_user_id;
DROP TABLE IF EXISTS group_members;
-- +goose StatementEnd
