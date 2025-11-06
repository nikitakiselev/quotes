package repository

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"quotes-backend/internal/models"
)

// QuoteRepository определяет интерфейс для работы с цитатами
type QuoteRepository interface {
	GetRandom() (*models.Quote, error)
	GetAll(page, pageSize int, search string) ([]models.Quote, int, error)
	GetByID(id string) (*models.Quote, error)
	Create(quote *models.Quote) error
	Update(id string, quote *models.Quote) error
	Delete(id string) error
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
		SELECT id, text, author, created_at, updated_at 
		FROM quotes 
		ORDER BY RANDOM() 
		LIMIT 1
	`

	var quote models.Quote
	err := r.db.QueryRow(query).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
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
		SELECT id, text, author, created_at, updated_at 
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
		SELECT id, text, author, created_at, updated_at 
		FROM quotes 
		WHERE id = $1
	`

	var quote models.Quote
	err := r.db.QueryRow(query, id).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
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
		INSERT INTO quotes (id, text, author, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	quote.CreatedAt = now
	quote.UpdatedAt = now

	_, err := r.db.Exec(
		query,
		quote.ID,
		quote.Text,
		quote.Author,
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

// CalculateTotalPages вычисляет общее количество страниц
func CalculateTotalPages(total, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

