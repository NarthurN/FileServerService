-- +goose Up
CREATE TABLE tokens (
    id VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4()::text,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,  -- Изменено с value на token
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true  -- Добавлено отсутствующее поле
);

CREATE INDEX idx_tokens_token ON tokens(token);  -- Изменено имя индекса
CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);
CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_is_active ON tokens(is_active);  -- Добавлен индекс для is_active

-- +goose Down
DROP TABLE IF EXISTS tokens;
