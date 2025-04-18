-- +goose Up
CREATE TABLE IF NOT EXISTS auth (
    id SERIAL PRIMARY KEY,
    guid VARCHAR(36) NOT NULL UNIQUE, -- Уникальный идентификатор пользователя
    hashed_token BYTEA NOT NULL, -- Хешированный Refresh Token (bcrypt)
    created_at TIMESTAMP NOT NULL DEFAULT now() -- Время создания записи
);

-- +goose Down
DROP TABLE IF EXISTS auth;