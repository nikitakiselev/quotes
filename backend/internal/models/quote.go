package models

import (
	"time"

	"github.com/google/uuid"
)

// Quote представляет цитату в системе
type Quote struct {
	ID        string    `json:"id" db:"id"`
	Text      string    `json:"text" db:"text"`
	Author    string    `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateQuoteRequest представляет запрос на создание цитаты
type CreateQuoteRequest struct {
	Text   string `json:"text" binding:"required"`
	Author string `json:"author" binding:"required"`
}

// UpdateQuoteRequest представляет запрос на обновление цитаты
type UpdateQuoteRequest struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

// QuoteResponse представляет ответ API с цитатой
type QuoteResponse struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginatedQuotesResponse представляет ответ API с пагинацией
type PaginatedQuotesResponse struct {
	Quotes      []QuoteResponse `json:"quotes"`
	Total       int             `json:"total"`
	Page        int             `json:"page"`
	PageSize    int             `json:"page_size"`
	TotalPages  int             `json:"total_pages"`
}

// ToResponse преобразует Quote в QuoteResponse
func (q *Quote) ToResponse() QuoteResponse {
	return QuoteResponse{
		ID:        q.ID,
		Text:      q.Text,
		Author:    q.Author,
		CreatedAt: q.CreatedAt,
		UpdatedAt: q.UpdatedAt,
	}
}

// NewQuote создает новую цитату из запроса
func NewQuote(req CreateQuoteRequest) *Quote {
	return &Quote{
		ID:     uuid.New().String(),
		Text:   req.Text,
		Author: req.Author,
	}
}

