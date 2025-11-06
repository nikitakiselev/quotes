package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"quotes-backend/internal/models"
	"quotes-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// QuoteHandler обрабатывает HTTP запросы для цитат
type QuoteHandler struct {
	repo repository.QuoteRepository
}

// NewQuoteHandler создает новый экземпляр обработчика
func NewQuoteHandler(repo repository.QuoteRepository) *QuoteHandler {
	return &QuoteHandler{repo: repo}
}

// getUserIP получает IP адрес пользователя из запроса
func getUserIP(c *gin.Context) string {
	userIP := c.ClientIP()
	// Если за прокси, пытаемся получить реальный IP
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		// Берем первый IP из списка
		if ips := strings.Split(forwarded, ","); len(ips) > 0 {
			userIP = strings.TrimSpace(ips[0])
		}
	}
	return userIP
}

// GetRandom возвращает случайную цитату
// @Summary Получить случайную цитату
// @Description Возвращает одну случайную цитату из базы данных
// @Tags quotes
// @Accept json
// @Produce json
// @Success 200 {object} models.QuoteResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/random [get]
func (h *QuoteHandler) GetRandom(c *gin.Context) {
	quote, err := h.repo.GetRandom()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, лайкнул ли текущий пользователь эту цитату
	userIP := getUserIP(c)
	isLiked, _ := h.repo.IsLiked(quote.ID, userIP)

	c.JSON(http.StatusOK, quote.ToResponse(isLiked))
}

// GetAll возвращает все цитаты с пагинацией
// @Summary Получить все цитаты
// @Description Возвращает список цитат с пагинацией и возможностью поиска
// @Tags quotes
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(10)
// @Param search query string false "Поисковый запрос"
// @Success 200 {object} models.PaginatedQuotesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes [get]
func (h *QuoteHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	quotes, total, err := h.repo.GetAll(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверяем статус лайка для каждой цитаты
	userIP := getUserIP(c)
	responses := make([]models.QuoteResponse, len(quotes))
	for i, quote := range quotes {
		isLiked, _ := h.repo.IsLiked(quote.ID, userIP)
		responses[i] = quote.ToResponse(isLiked)
	}

	totalPages := repository.CalculateTotalPages(total, pageSize)

	c.JSON(http.StatusOK, models.PaginatedQuotesResponse{
		Quotes:     responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetByID возвращает цитату по ID
// @Summary Получить цитату по ID
// @Description Возвращает цитату с указанным ID
// @Tags quotes
// @Accept json
// @Produce json
// @Param id path string true "ID цитаты"
// @Success 200 {object} models.QuoteResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id} [get]
func (h *QuoteHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	quote, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, лайкнул ли текущий пользователь эту цитату
	userIP := getUserIP(c)
	isLiked, _ := h.repo.IsLiked(quote.ID, userIP)

	c.JSON(http.StatusOK, quote.ToResponse(isLiked))
}

// Create создает новую цитату
// @Summary Создать новую цитату
// @Description Создает новую цитату в базе данных
// @Tags quotes
// @Accept json
// @Produce json
// @Param quote body models.CreateQuoteRequest true "Данные цитаты"
// @Success 201 {object} models.QuoteResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes [post]
func (h *QuoteHandler) Create(c *gin.Context) {
	var req models.CreateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quote := models.NewQuote(req)
	if err := h.repo.Create(quote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Новая цитата не может быть лайкнута
	c.JSON(http.StatusCreated, quote.ToResponse(false))
}

// Update обновляет существующую цитату
// @Summary Обновить цитату
// @Description Обновляет существующую цитату по ID
// @Tags quotes
// @Accept json
// @Produce json
// @Param id path string true "ID цитаты"
// @Param quote body models.UpdateQuoteRequest true "Обновленные данные цитаты"
// @Success 200 {object} models.QuoteResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id} [put]
func (h *QuoteHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем существующую цитату
	quote, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Обновляем поля
	if req.Text != "" {
		quote.Text = req.Text
	}
	if req.Author != "" {
		quote.Author = req.Author
	}

	if err := h.repo.Update(id, quote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем обновленную цитату
	updatedQuote, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, лайкнул ли текущий пользователь эту цитату
	userIP := getUserIP(c)
	isLiked, _ := h.repo.IsLiked(updatedQuote.ID, userIP)

	c.JSON(http.StatusOK, updatedQuote.ToResponse(isLiked))
}

// Delete удаляет цитату
// @Summary Удалить цитату
// @Description Удаляет цитату по ID
// @Tags quotes
// @Accept json
// @Produce json
// @Param id path string true "ID цитаты"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id} [delete]
func (h *QuoteHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Like ставит лайк цитате
// @Summary Поставить лайк цитате
// @Description Увеличивает количество лайков у цитаты на 1. Предотвращает множественные лайки от одного пользователя
// @Tags quotes
// @Accept json
// @Produce json
// @Param id path string true "ID цитаты"
// @Success 200 {object} models.QuoteResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id}/like [put]
func (h *QuoteHandler) Like(c *gin.Context) {
	id := c.Param("id")

	// Получаем IP адрес пользователя
	userIP := getUserIP(c)
	userAgent := c.GetHeader("User-Agent")

	if err := h.repo.Like(id, userIP, userAgent); err != nil {
		// Проверяем, это ошибка "уже лайкнуто" или другая
		if strings.Contains(err.Error(), "already liked") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Вы уже поставили лайк этой цитате"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Получаем обновленную цитату
	quote, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// После лайка пользователь точно лайкнул эту цитату
	c.JSON(http.StatusOK, quote.ToResponse(true))
}

// GetTopWeekly возвращает топ цитату за неделю
// @Summary Получить топ цитату за неделю
// @Description Возвращает цитату с наибольшим количеством лайков за последние 7 дней
// @Tags quotes
// @Accept json
// @Produce json
// @Success 200 {object} models.QuoteResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/top/weekly [get]
func (h *QuoteHandler) GetTopWeekly(c *gin.Context) {
	quote, err := h.repo.GetTopWeekly()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, лайкнул ли текущий пользователь эту цитату
	userIP := getUserIP(c)
	isLiked, _ := h.repo.IsLiked(quote.ID, userIP)

	c.JSON(http.StatusOK, quote.ToResponse(isLiked))
}

// GetTopAllTime возвращает топ цитату за всё время
// @Summary Получить топ цитату за всё время
// @Description Возвращает цитату с наибольшим количеством лайков за всё время
// @Tags quotes
// @Accept json
// @Produce json
// @Success 200 {object} models.QuoteResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/top/alltime [get]
func (h *QuoteHandler) GetTopAllTime(c *gin.Context) {
	quote, err := h.repo.GetTopAllTime()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, лайкнул ли текущий пользователь эту цитату
	userIP := getUserIP(c)
	isLiked, _ := h.repo.IsLiked(quote.ID, userIP)

	c.JSON(http.StatusOK, quote.ToResponse(isLiked))
}

// IsLiked проверяет, лайкнул ли текущий пользователь цитату
// @Summary Проверить статус лайка
// @Description Проверяет, поставил ли текущий пользователь лайк указанной цитате
// @Tags quotes
// @Accept json
// @Produce json
// @Param id path string true "ID цитаты"
// @Success 200 {object} map[string]bool
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id}/is-liked [get]
func (h *QuoteHandler) IsLiked(c *gin.Context) {
	id := c.Param("id")

	// Получаем IP адрес пользователя
	userIP := getUserIP(c)

	isLiked, err := h.repo.IsLiked(id, userIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_liked": isLiked})
}

// ResetLikes сбрасывает все лайки
// @Summary Сбросить все лайки
// @Description Обнуляет счетчики лайков у всех цитат и удаляет все записи о лайках
// @Tags quotes
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/likes/reset [delete]
func (h *QuoteHandler) ResetLikes(c *gin.Context) {
	if err := h.repo.ResetLikes(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Все лайки успешно сброшены"})
}

