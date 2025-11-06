-- Инициализация базы данных
-- Этот файл выполняется автоматически при первом запуске контейнера PostgreSQL

-- Создание таблицы quotes
CREATE TABLE IF NOT EXISTS quotes (
    id VARCHAR(36) PRIMARY KEY,
    text TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_quotes_text ON quotes USING gin(to_tsvector('russian', text));
CREATE INDEX IF NOT EXISTS idx_quotes_author ON quotes(author);
CREATE INDEX IF NOT EXISTS idx_quotes_created_at ON quotes(created_at DESC);

-- Вставка тестовых данных
INSERT INTO quotes (id, text, author, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Единственный способ делать великую работу — это любить то, что ты делаешь.', 'Стив Джобс', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440001', 'Инновация отличает лидера от последователя.', 'Стив Джобс', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440002', 'Ваше время ограничено, не тратьте его, живя чужой жизнью.', 'Стив Джобс', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440003', 'Будьте голодными. Будьте безрассудными.', 'Стив Джобс', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440004', 'Простота — это высшая форма изысканности.', 'Леонардо да Винчи', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440005', 'Жизнь — это то, что происходит с тобой, пока ты строишь планы.', 'Джон Леннон', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440006', 'Успех — это способность идти от неудачи к неудаче, не теряя энтузиазма.', 'Уинстон Черчилль', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440007', 'Будущее принадлежит тем, кто верит в красоту своих мечтаний.', 'Элеонора Рузвельт', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440008', 'Единственный человек, которым вы должны стать — это тот, кем вы решили стать.', 'Ральф Уолдо Эмерсон', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440009', 'Не важно, как медленно ты идешь, до тех пор, пока ты не останавливаешься.', 'Конфуций', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440010', 'Лучшее время посадить дерево было 20 лет назад. Следующее лучшее время — сейчас.', 'Китайская мудрость', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440011', 'Два самых важных дня в твоей жизни: день, когда ты родился, и день, когда ты понял зачем.', 'Марк Твен', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440012', 'Качество — это не действие, это привычка.', 'Аристотель', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440013', 'Стремитесь не к успеху, а к ценностям, которые он дает.', 'Альберт Эйнштейн', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440014', 'Единственный способ иметь друга — быть им.', 'Ральф Уолдо Эмерсон', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

