-- +goose Up
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(255),
    file_path VARCHAR(500),
    is_file BOOLEAN DEFAULT true,
    is_public BOOLEAN DEFAULT false,
    json_data JSONB,
    grants JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для производительности (критично для задания!)
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_created_at ON documents(created_at);
CREATE INDEX idx_documents_name ON documents(name);
CREATE INDEX idx_documents_is_public ON documents(is_public);
CREATE INDEX idx_documents_grants ON documents USING GIN(grants);

-- +goose Down
DROP TABLE IF EXISTS documents;
