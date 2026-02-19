-- +goose Up
-- +goose StatementBegin
CREATE TABLE whatsapp_whitelisted_groups
(
  jid        TEXT NOT NULL,
  server_jid TEXT NOT NULL,
  PRIMARY KEY (jid, server_jid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS whatsapp_whitelisted_groups;
-- +goose StatementEnd

