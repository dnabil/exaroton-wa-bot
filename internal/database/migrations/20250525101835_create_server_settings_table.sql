-- +goose Up
-- +goose StatementBegin
CREATE TABLE server_settings (
    key TEXT PRIMARY KEY,
	value TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS server_settings;
-- +goose StatementEnd
