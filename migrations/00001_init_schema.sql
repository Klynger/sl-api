-- +goose Up
-- +goose StatementBegin
CREATE TABLE products
(
  id UUID PRIMARY KEY,
  product_name TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_products_deleted_at ON products (deleted_at);

CREATE TABLE users
(
  id UUID PRIMARY KEY,
  user_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  username VARCHAR(50) NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP,
  CONSTRAINT users_username_unique UNIQUE (username)
);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);

CREATE TABLE groups
(
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL CHECK (char_length(name) <= 255),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_groups_deleted_at ON groups (deleted_at);

CREATE TABLE group_members
(
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  group_id UUID NOT NULL,
  roles TEXT [] DEFAULT ARRAY['member']::TEXT [] NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
  CONSTRAINT unique_user_group UNIQUE (user_id, group_id)
);

CREATE INDEX idx_group_members_user_id ON group_members (user_id);
CREATE INDEX idx_group_members_group_id ON group_members (group_id);
CREATE INDEX idx_group_members_deleted_at ON group_members (deleted_at);

CREATE TABLE invites
(
  id UUID PRIMARY KEY,
  group_id UUID NOT NULL,
  sender_id UUID NOT NULL,
  invited_user_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP,
  CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
  CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_invited_user FOREIGN KEY (invited_user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT unique_invite UNIQUE (group_id, invited_user_id)
);

CREATE INDEX idx_invites_group_id ON invites (group_id);
CREATE INDEX idx_invites_sender_id ON invites (sender_id);
CREATE INDEX idx_invites_invited_user_id ON invites (invited_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invites;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
