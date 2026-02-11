-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT users_username_unique;
-- +goose StatementEnd
