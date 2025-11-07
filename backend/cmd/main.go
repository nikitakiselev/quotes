package main

import (
	"log"
	"os"

	"quotes-backend/internal/config"
	"quotes-backend/internal/database"
	"quotes-backend/internal/handlers"
	"quotes-backend/internal/repository"
	"quotes-backend/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Устанавливаем production режим для Gin (отключает debug режим)
	gin.SetMode(gin.ReleaseMode)

	// Инициализация конфигурации
	cfg := config.Load()

	// Инициализация базы данных
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Выполнение миграций
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Инициализация репозитория
	quoteRepo := repository.NewQuoteRepository(db)

	// Инициализация обработчиков
	quoteHandler := handlers.NewQuoteHandler(quoteRepo)

	// Настройка роутера
	r := router.SetupRouter(quoteHandler, cfg)

	// Запуск сервера
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

