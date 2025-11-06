package repository

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"quotes-backend/internal/models"

	"github.com/google/uuid"
)

// QuoteRepository определяет интерфейс для работы с цитатами
type QuoteRepository interface {
	GetRandom() (*models.Quote, error)
	GetAll(page, pageSize int, search string) ([]models.Quote, int, error)
	GetByID(id string) (*models.Quote, error)
	Create(quote *models.Quote) error
	Update(id string, quote *models.Quote) error
	Delete(id string) error
	Like(id string, userIP, userAgent string) error
	IsLiked(id string, userIP string) (bool, error)
	GetTopWeekly() (*models.Quote, error)
	GetTopAllTime() (*models.Quote, error)
	ResetLikes() error
}

type quoteRepository struct {
	db *sql.DB
}

// NewQuoteRepository создает новый экземпляр репозитория
func NewQuoteRepository(db *sql.DB) QuoteRepository {
	return &quoteRepository{db: db}
}

// GetRandom возвращает случайную цитату
func (r *quoteRepository) GetRandom() (*models.Quote, error) {
	query := `
		SELECT id, text, author, likes_count, created_at, updated_at 
		FROM quotes 
		ORDER BY RANDOM() 
		LIMIT 1
	`

	var quote models.Quote
	err := r.db.QueryRow(query).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.LikesCount,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no quotes found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return &quote, nil
}

// GetAll возвращает все цитаты с пагинацией и поиском
func (r *quoteRepository) GetAll(page, pageSize int, search string) ([]models.Quote, int, error) {
	// Подсчет общего количества
	var total int
	countQuery := "SELECT COUNT(*) FROM quotes"
	args := []interface{}{}

	if search != "" {
		countQuery += " WHERE text ILIKE $1 OR author ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count quotes: %w", err)
	}

	// Получение цитат с пагинацией
	offset := (page - 1) * pageSize
	query := `
		SELECT id, text, author, likes_count, created_at, updated_at 
		FROM quotes
	`

	if search != "" {
		query += " WHERE text ILIKE $1 OR author ILIKE $1"
		query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, pageSize, offset)
	} else {
		query += " ORDER BY created_at DESC LIMIT $1 OFFSET $2"
		args = []interface{}{pageSize, offset}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get quotes: %w", err)
	}
	defer rows.Close()

	var quotes []models.Quote
	for rows.Next() {
		var quote models.Quote
		if err := rows.Scan(
			&quote.ID,
			&quote.Text,
			&quote.Author,
			&quote.LikesCount,
			&quote.CreatedAt,
			&quote.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, quote)
	}

	return quotes, total, nil
}

// GetByID возвращает цитату по ID
func (r *quoteRepository) GetByID(id string) (*models.Quote, error) {
	query := `
		SELECT id, text, author, likes_count, created_at, updated_at 
		FROM quotes 
		WHERE id = $1
	`

	var quote models.Quote
	err := r.db.QueryRow(query, id).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.LikesCount,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("quote not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	return &quote, nil
}

// Create создает новую цитату
func (r *quoteRepository) Create(quote *models.Quote) error {
	query := `
		INSERT INTO quotes (id, text, author, likes_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	quote.CreatedAt = now
	quote.UpdatedAt = now
	quote.LikesCount = 0

	_, err := r.db.Exec(
		query,
		quote.ID,
		quote.Text,
		quote.Author,
		quote.LikesCount,
		quote.CreatedAt,
		quote.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create quote: %w", err)
	}

	return nil
}

// Update обновляет существующую цитату
func (r *quoteRepository) Update(id string, quote *models.Quote) error {
	query := `
		UPDATE quotes 
		SET text = $1, author = $2, updated_at = $3
		WHERE id = $4
	`

	quote.UpdatedAt = time.Now()

	result, err := r.db.Exec(query, quote.Text, quote.Author, quote.UpdatedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update quote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quote not found")
	}

	return nil
}

// Delete удаляет цитату
func (r *quoteRepository) Delete(id string) error {
	query := `DELETE FROM quotes WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete quote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quote not found")
	}

	return nil
}

// Like увеличивает количество лайков у цитаты
// userIP и userAgent используются для предотвращения множественных лайков
// Защита от накрутки работает на нескольких уровнях:
// 1. Проверка существующего лайка перед транзакцией
// 2. Уникальное ограничение UNIQUE(quote_id, user_ip) в таблице likes
// 3. ON CONFLICT DO NOTHING в INSERT для предотвращения дубликатов
// 4. Транзакция обеспечивает атомарность операции
// Это защищает от накрутки даже при прямых HTTP запросах
func (r *quoteRepository) Like(id string, userIP, userAgent string) error {
	// Начинаем транзакцию сразу для изоляции
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Проверяем, не лайкал ли уже этот пользователь эту цитату (внутри транзакции)
	checkQuery := `
		SELECT id FROM likes 
		WHERE quote_id = $1 AND user_ip = $2
		FOR UPDATE
	`
	var existingLikeID string
	err = tx.QueryRow(checkQuery, id, userIP).Scan(&existingLikeID)
	
	if err == nil {
		// Лайк уже существует, возвращаем ошибку
		return fmt.Errorf("you have already liked this quote")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing like: %w", err)
	}

	// Увеличиваем счетчик лайков
	updateQuery := `
		UPDATE quotes 
		SET likes_count = likes_count + 1, updated_at = $1
		WHERE id = $2
	`
	result, err := tx.Exec(updateQuery, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update likes count: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quote not found")
	}

	// Сохраняем информацию о лайке
	// ON CONFLICT DO NOTHING - дополнительная защита от race condition
	likeID := uuid.New().String()
	insertQuery := `
		INSERT INTO likes (id, quote_id, user_ip, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (quote_id, user_ip) DO NOTHING
	`
	_, err = tx.Exec(insertQuery, likeID, id, userIP, userAgent, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save like: %w", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// IsLiked проверяет, лайкнул ли пользователь цитату
func (r *quoteRepository) IsLiked(id string, userIP string) (bool, error) {
	query := `SELECT COUNT(*) FROM likes WHERE quote_id = $1 AND user_ip = $2`
	var count int
	err := r.db.QueryRow(query, id, userIP).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}
	return count > 0, nil
}

// GetTopWeekly возвращает цитату с наибольшим количеством лайков за последнюю неделю
func (r *quoteRepository) GetTopWeekly() (*models.Quote, error) {
	query := `
		SELECT id, text, author, likes_count, created_at, updated_at 
		FROM quotes 
		WHERE created_at >= NOW() - INTERVAL '7 days'
		ORDER BY likes_count DESC, created_at DESC
		LIMIT 1
	`

	var quote models.Quote
	err := r.db.QueryRow(query).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.LikesCount,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no quotes found for the last week")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get top weekly quote: %w", err)
	}

	return &quote, nil
}

// GetTopAllTime возвращает цитату с наибольшим количеством лайков за всё время
func (r *quoteRepository) GetTopAllTime() (*models.Quote, error) {
	query := `
		SELECT id, text, author, likes_count, created_at, updated_at 
		FROM quotes 
		ORDER BY likes_count DESC, created_at DESC
		LIMIT 1
	`

	var quote models.Quote
	err := r.db.QueryRow(query).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.LikesCount,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no quotes found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get top all time quote: %w", err)
	}

	return &quote, nil
}

// ResetLikes сбрасывает все лайки: обнуляет счетчики и удаляет записи из таблицы likes
func (r *quoteRepository) ResetLikes() error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Обнуляем счетчики лайков у всех цитат
	updateQuery := `
		UPDATE quotes 
		SET likes_count = 0, updated_at = $1
	`
	_, err = tx.Exec(updateQuery, time.Now())
	if err != nil {
		return fmt.Errorf("failed to reset likes count: %w", err)
	}

	// Удаляем все записи из таблицы likes
	deleteQuery := `DELETE FROM likes`
	_, err = tx.Exec(deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to delete likes: %w", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CalculateTotalPages вычисляет общее количество страниц
func CalculateTotalPages(total, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

