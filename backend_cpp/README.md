# C++ Backend для Quotes API

Высокопроизводительный backend на C++ с использованием фреймворка Drogon.

## Особенности

- **Drogon** - быстрый async веб-фреймворк на C++17
- **libpq** - прямое подключение к PostgreSQL для максимальной производительности
- **Оптимизация компиляции**: `-O3 -march=native -flto`
- **Batch запросы** для устранения N+1 проблемы
- **Prepared statements** для всех SQL запросов
- **Автоматическое определение количества потоков** (по количеству CPU ядер)

## Сборка

```bash
docker compose build backend_cpp
```

## Запуск

```bash
docker compose up -d backend_cpp
```

Backend будет доступен на порту 8083.

## Endpoints

- `GET /health` - health check
- `GET /api/quotes/random` - случайная цитата
- `GET /api/quotes` - все цитаты с пагинацией
- `GET /api/quotes/:id` - цитата по ID
- `POST /api/quotes` - создать цитату
- `PUT /api/quotes/:id` - обновить цитату
- `DELETE /api/quotes/:id` - удалить цитату
- `PUT /api/quotes/:id/like` - поставить лайк
- `GET /api/quotes/top/weekly` - топ за неделю
- `GET /api/quotes/top/alltime` - топ за всё время
- `DELETE /api/quotes/likes/reset` - сбросить все лайки

## Заголовки

Все ответы содержат заголовок `X-Backend: cpp`.

