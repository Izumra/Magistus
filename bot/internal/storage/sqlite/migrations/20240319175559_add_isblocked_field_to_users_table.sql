-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN is_blocked BOOL NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_blocked;
-- +goose StatementEnd
