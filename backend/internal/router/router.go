package router

import (
	"net/http"
	"os"
	"path/filepath"
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
			quotes.GET("/random", quoteHandler.GetRandom)
			quotes.GET("", quoteHandler.GetAll)
			quotes.GET("/:id", quoteHandler.GetByID)
			quotes.POST("", quoteHandler.Create)
			quotes.PUT("/:id", quoteHandler.Update)
			quotes.DELETE("/:id", quoteHandler.Delete)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Отдача статических файлов фронтенда
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./frontend/dist"
	}

	// Проверяем существование директории со статикой
	if _, err := os.Stat(staticDir); err == nil {
		// Статические файлы из папки assets (JS, CSS и т.д.)
		assetsPath := filepath.Join(staticDir, "assets")
		if _, err := os.Stat(assetsPath); err == nil {
			r.Static("/assets", assetsPath)
		}

		// Отдельные статические файлы в корне
		faviconPath := filepath.Join(staticDir, "favicon.ico")
		if _, err := os.Stat(faviconPath); err == nil {
			r.StaticFile("/favicon.ico", faviconPath)
		}

		viteSvgPath := filepath.Join(staticDir, "vite.svg")
		if _, err := os.Stat(viteSvgPath); err == nil {
			r.StaticFile("/vite.svg", viteSvgPath)
		}

		// SPA routing - все остальные запросы отдаем index.html
		r.NoRoute(func(c *gin.Context) {
			// Если запрос к API, возвращаем 404
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}

			// Для всех остальных запросов отдаем index.html (SPA routing)
			indexPath := filepath.Join(staticDir, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
			} else {
				c.String(http.StatusNotFound, "Frontend not found")
			}
		})
	}

	return r
}

