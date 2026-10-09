-- +goose Up
-- +goose StatementBegin
-- A list_item is ephemeral: "this product is currently on this list". Its
-- history has no value, so it moves from soft delete to hard delete. Removing
-- an item deletes the row; re-adding inserts a fresh one. The unique
-- (list_id, product_id) constraint now simply prevents a product appearing
-- twice on a list, with no soft-deleted rows to span.
DROP INDEX IF EXISTS idx_list_items_deleted_at;
ALTER TABLE list_items DROP COLUMN deleted_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE list_items ADD COLUMN deleted_at TIMESTAMP;
CREATE INDEX idx_list_items_deleted_at ON list_items (deleted_at);
-- +goose StatementEnd
