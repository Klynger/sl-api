-- +goose Up
-- +goose StatementBegin
ALTER TABLE products RENAME COLUMN name TO product_name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products RENAME COLUMN product_name TO name;
-- +goose StatementEnd
