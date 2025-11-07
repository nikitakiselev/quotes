package router

import (
	"net/http"
	"strings"

	"quotes-backend/internal/config"
	"quotes-backend/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает и возвращает роутер
func SetupRouter(quoteHandler *handlers.QuoteHandler, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Настройка CORS
	corsConfig := cors.DefaultConfig()
	if cfg.CORSOrigin == "*" {
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowWildcard = true
	} else {
		corsConfig.AllowOrigins = []string{cfg.CORSOrigin}
		corsConfig.AllowCredentials = true
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	corsConfig.AllowBrowserExtensions = true

	r.Use(cors.New(corsConfig))

	// API routes
	api := r.Group("/api")
	{
		quotes := api.Group("/quotes")
		{
			// Специфичные роуты должны быть раньше параметризованных
			quotes.GET("/random", quoteHandler.GetRandom)
			quotes.GET("/top/weekly", quoteHandler.GetTopWeekly)
			quotes.GET("/top/alltime", quoteHandler.GetTopAllTime)
			quotes.DELETE("/likes/reset", quoteHandler.ResetLikes)
			quotes.GET("", quoteHandler.GetAll)
			quotes.POST("", quoteHandler.Create)
			// Параметризованные роуты в конце
			quotes.PUT("/:id/like", quoteHandler.Like)
			quotes.GET("/:id", quoteHandler.GetByID)
			quotes.PUT("/:id", quoteHandler.Update)
			quotes.DELETE("/:id", quoteHandler.Delete)
		}
	}


	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API-only backend - не отдаем статику
	r.NoRoute(func(c *gin.Context) {
		// Если запрос к API, возвращаем 404
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		// Для всех остальных запросов возвращаем 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})

	return r
}

