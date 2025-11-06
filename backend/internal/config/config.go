package config

import (
	"os"
)

// Config содержит конфигурацию приложения
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	APIPort    string
	CORSOrigin string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "quotes_user"),
		DBPassword: getEnv("DB_PASSWORD", "quotes_password"),
		DBName:     getEnv("DB_NAME", "quotes_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		APIPort:    getEnv("API_PORT", "8080"),
		CORSOrigin: getEnv("CORS_ORIGIN", "http://localhost:3000"),
	}
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

