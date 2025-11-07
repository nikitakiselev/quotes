package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"quotes-backend/internal/config"

	_ "github.com/lib/pq"
)

// Connect устанавливает соединение с базой данных PostgreSQL
func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Оптимизация пула соединений для высокой производительности
	// MaxOpenConns - максимальное количество открытых соединений
	// Для высоконагруженного сервера устанавливаем 25 (по умолчанию 0 = неограниченно, но это плохо)
	db.SetMaxOpenConns(25)
	
	// MaxIdleConns - максимальное количество неактивных соединений в пуле
	// Должно быть меньше MaxOpenConns
	db.SetMaxIdleConns(10)
	
	// ConnMaxLifetime - максимальное время жизни соединения
	// Переподключаемся каждые 5 минут для предотвращения проблем с таймаутами
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// ConnMaxIdleTime - максимальное время простоя соединения
	// Закрываем неиспользуемые соединения через 10 минут
	db.SetConnMaxIdleTime(10 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations выполняет миграции из директории migrations
func RunMigrations(db *sql.DB) error {
	// В Docker контейнере миграции монтируются в /app/db/migrations
	// В локальной разработке используем относительный путь
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		// Пробуем разные пути
		possiblePaths := []string{
			"/app/db/migrations",           // Docker
			"../../db/migrations",          // Локально из cmd/
			"../db/migrations",             // Альтернативный локальный путь
		}
		
		var absPath string
		var err error
		for _, path := range possiblePaths {
			absPath, err = filepath.Abs(path)
			if err == nil {
				if _, err := os.Stat(absPath); err == nil {
					migrationsDir = absPath
					break
				}
			}
		}
		
		if migrationsDir == "" {
			return fmt.Errorf("migrations directory not found")
		}
	} else {
		var err error
		migrationsDir, err = filepath.Abs(migrationsDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
	}

	// Читаем файлы миграций
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	if len(files) == 0 {
		// Если миграций нет, это не ошибка - возможно они уже применены через init скрипт
		return nil
	}

	// Сортируем файлы по имени
	for _, file := range files {
		migrationSQL, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(migrationSQL)); err != nil {
			// Игнорируем ошибки "уже существует" для таблиц и индексов
			if !strings.Contains(strings.ToLower(err.Error()), "already exists") {
				return fmt.Errorf("failed to execute migration %s: %w", file, err)
			}
		}
	}

	return nil
}

