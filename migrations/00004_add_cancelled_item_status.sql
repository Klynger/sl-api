-- +goose Up
-- +goose StatementBegin
-- Add a "cancelled" item status: the shopper gave up on getting this product but
-- left it on the list, which is distinct from deleting the row off the list.
-- Statuses are an extensible enum, so widening the allowed set is a simple swap
-- of the CHECK constraint.
ALTER TABLE list_items DROP CONSTRAINT list_items_status_check;
ALTER TABLE list_items ADD CONSTRAINT list_items_status_check
  CHECK (status IN ('pending', 'in_cart', 'purchased', 'cancelled'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE list_items SET status = 'pending' WHERE status = 'cancelled';
ALTER TABLE list_items DROP CONSTRAINT list_items_status_check;
ALTER TABLE list_items ADD CONSTRAINT list_items_status_check
  CHECK (status IN ('pending', 'in_cart', 'purchased'));
-- +goose StatementEnd
