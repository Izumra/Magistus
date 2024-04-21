-- +goose Up
-- +goose StatementBegin
ALTER TABLE charts ADD COLUMN interpritation TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE charts DROP COLUMN interpritation;
-- +goose StatementEnd
