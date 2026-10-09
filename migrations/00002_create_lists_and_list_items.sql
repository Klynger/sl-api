-- +goose Up
-- +goose StatementBegin
ALTER TABLE products ADD COLUMN predefined_unit VARCHAR(50);

CREATE TABLE lists
(
  id UUID PRIMARY KEY,
  group_id UUID NOT NULL,
  created_by UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  status TEXT NOT NULL DEFAULT 'planning',
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP,
  CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
  CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT lists_status_check CHECK (status IN ('planning', 'active', 'completed'))
);

CREATE INDEX idx_lists_group_id ON lists (group_id);
CREATE INDEX idx_lists_deleted_at ON lists (deleted_at);

-- The unique constraint below is deliberately NOT partial on deleted_at: a
-- (list_id, product_id) pair is unique for all time, including soft-deleted
-- rows, so re-adding a removed product revives its existing row instead of
-- inserting a duplicate.
CREATE TABLE list_items
(
  id UUID PRIMARY KEY,
  list_id UUID NOT NULL,
  product_id UUID NOT NULL,
  added_by UUID NOT NULL,
  quantity INTEGER NOT NULL DEFAULT 1,
  unit VARCHAR(50),
  status TEXT NOT NULL DEFAULT 'pending',
  note TEXT,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP,
  CONSTRAINT fk_list FOREIGN KEY (list_id) REFERENCES lists (id) ON DELETE CASCADE,
  CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE RESTRICT,
  CONSTRAINT fk_added_by FOREIGN KEY (added_by) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT unique_list_product UNIQUE (list_id, product_id),
  CONSTRAINT list_items_status_check CHECK (status IN ('pending', 'in_cart', 'purchased')),
  CONSTRAINT list_items_quantity_check CHECK (quantity > 0)
);

CREATE INDEX idx_list_items_list_id ON list_items (list_id);
CREATE INDEX idx_list_items_product_id ON list_items (product_id);
CREATE INDEX idx_list_items_deleted_at ON list_items (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS list_items;
DROP TABLE IF EXISTS lists;
ALTER TABLE products DROP COLUMN IF EXISTS predefined_unit;
-- +goose StatementEnd
