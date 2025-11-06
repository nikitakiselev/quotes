.PHONY: help build up down restart logs clean

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Собрать Docker образы
	docker-compose build

up: ## Запустить все сервисы
	docker-compose up -d

down: ## Остановить все сервисы
	docker-compose down

restart: ## Перезапустить все сервисы
	docker-compose restart

logs: ## Показать логи всех сервисов
	docker-compose logs -f

logs-backend: ## Показать логи бэкенда (включает фронтенд)
	docker-compose logs -f backend

logs-admin: ## Показать логи админки
	docker-compose logs -f admin

logs-db: ## Показать логи базы данных
	docker-compose logs -f postgres

clean: ## Остановить и удалить все контейнеры, volumes и сети
	docker-compose down -v

ps: ## Показать статус контейнеров
	docker-compose ps

db-shell: ## Подключиться к базе данных
	docker-compose exec postgres psql -U quotes_user -d quotes_db

