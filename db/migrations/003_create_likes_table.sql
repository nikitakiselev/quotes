-- Создание таблицы для отслеживания лайков
CREATE TABLE IF NOT EXISTS likes (
    id VARCHAR(36) PRIMARY KEY,
    quote_id VARCHAR(36) NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    user_ip VARCHAR(45) NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(quote_id, user_ip)
);

-- Создание индекса для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_likes_quote_id ON likes(quote_id);
CREATE INDEX IF NOT EXISTS idx_likes_user_ip ON likes(user_ip);

