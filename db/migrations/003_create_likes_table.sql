-- Создание таблицы для отслеживания лайков
CREATE TABLE IF NOT EXISTS likes (
    id VARCHAR(36) PRIMARY KEY,
    quote_id VARCHAR(36) NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    user_ip VARCHAR(45) NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(quote_id, user_ip)
);

-- Если таблица уже существует, добавляем уникальное ограничение отдельно
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'likes_quote_id_user_ip_key'
    ) THEN
        ALTER TABLE likes ADD CONSTRAINT likes_quote_id_user_ip_key UNIQUE (quote_id, user_ip);
    END IF;
END $$;

-- Создание индекса для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_likes_quote_id ON likes(quote_id);
CREATE INDEX IF NOT EXISTS idx_likes_user_ip ON likes(user_ip);

