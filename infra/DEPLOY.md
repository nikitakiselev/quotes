# Инструкция по деплою на production сервер

## Первоначальная настройка

1. Подключитесь к серверу:
```bash
ssh nikitakiselev@192.168.88.40
```

2. Создайте директорию для проекта (если не существует):
```bash
mkdir -p /home/nikitakiselev/infra/services/quotes
```

3. Скопируйте скрипт деплоя на сервер или клонируйте репозиторий:
```bash
cd /home/nikitakiselev/infra/services/quotes
git clone git@github.com:nikitakiselev/quotes.git .
```

4. Создайте `.env` файл из примера:
```bash
cp .env.example .env
nano .env  # или используйте ваш любимый редактор
```

5. Отредактируйте `.env` файл, особенно важно:
   - Установить безопасный `ADMIN_PASSWORD`
   - Проверить настройки базы данных
   - Установить правильный `CORS_ORIGIN` для production (если нужно)
   - Установить правильные порты

6. Запустите деплой:
```bash
bash infra/deploy.sh
```

## Обновление проекта

Для обновления проекта после изменений в репозитории:

```bash
cd /home/nikitakiselev/infra/services/quotes
bash infra/deploy.sh
```

Скрипт автоматически:
- Обновит код из репозитория
- Пересоберет Docker образы
- Перезапустит контейнеры

## Ручной деплой

Если нужно выполнить деплой вручную:

```bash
cd /home/nikitakiselev/infra/services/quotes

# Обновление кода
git pull origin main

# Остановка контейнеров
docker-compose down

# Сборка образов
docker-compose build --no-cache

# Запуск контейнеров
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f
```

## Полезные команды

```bash
# Просмотр логов всех сервисов
docker-compose logs -f

# Просмотр логов конкретного сервиса
docker-compose logs -f backend
docker-compose logs -f admin
docker-compose logs -f postgres

# Остановка всех сервисов
docker-compose down

# Перезапуск сервисов
docker-compose restart

# Подключение к базе данных
docker-compose exec postgres psql -U quotes_user -d quotes_db

# Проверка использования ресурсов
docker stats
```

## Troubleshooting

### Проблемы с портами

Если порты заняты, измените их в `.env`:
```env
FRONTEND_PORT=3000
ADMIN_PORT=3001
DB_PORT=5432
```

### Проблемы с правами доступа

Убедитесь, что у пользователя есть права на Docker:
```bash
sudo usermod -aG docker nikitakiselev
# Перелогиньтесь после этого
```

### Проблемы с базой данных

Если база данных не запускается:
```bash
# Проверьте логи
docker-compose logs postgres

# Проверьте volumes
docker volume ls | grep quotes
```

### Очистка и пересборка

Если что-то пошло не так:
```bash
docker-compose down -v  # Удалит volumes тоже
docker-compose build --no-cache
docker-compose up -d
```

