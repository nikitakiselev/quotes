-- Добавление поля likes_count в таблицу quotes
ALTER TABLE quotes ADD COLUMN IF NOT EXISTS likes_count INTEGER DEFAULT 0 NOT NULL;

-- Создание индекса для быстрого поиска по лайкам
CREATE INDEX IF NOT EXISTS idx_quotes_likes_count ON quotes(likes_count DESC);

