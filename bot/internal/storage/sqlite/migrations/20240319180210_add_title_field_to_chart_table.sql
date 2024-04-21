-- +goose Up
-- +goose StatementBegin
ALTER TABLE charts ADD COLUMN title VARCHAR(30) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE charts DROP COLUMN title;
-- +goose StatementEnd
