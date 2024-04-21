-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS charts(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query TEXT NOT NULL,
    content TEXT DEFAULT '',
    id_creator BIGINT REFERENCES users(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS charts;
-- +goose StatementEnd
