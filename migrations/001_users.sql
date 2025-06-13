-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX users_chat_id_idx ON users (chat_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

