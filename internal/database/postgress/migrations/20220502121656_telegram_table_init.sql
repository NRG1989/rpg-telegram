-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS userservice.client_telegram
( id bigserial NOT NULL,
phone_number varchar(13), -- for Russian numbers only might be 12, but i left 13 for Ukrainians to test it too
chat_id bigint DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS userservice.client_telegram;
-- +goose StatementEnd
