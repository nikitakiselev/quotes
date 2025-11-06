-- Создание таблицы quotes
CREATE TABLE IF NOT EXISTS quotes (
    id VARCHAR(36) PRIMARY KEY,
    text TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создание индекса для поиска
CREATE INDEX IF NOT EXISTS idx_quotes_text ON quotes USING gin(to_tsvector('russian', text));
CREATE INDEX IF NOT EXISTS idx_quotes_author ON quotes(author);
CREATE INDEX IF NOT EXISTS idx_quotes_created_at ON quotes(created_at DESC);

