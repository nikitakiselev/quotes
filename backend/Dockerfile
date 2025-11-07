# Этап 1: Сборка бэкенда
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# Установка зависимостей для разработки
RUN apk add --no-cache git

# Копирование файлов зависимостей
COPY backend/go.mod ./
COPY backend/go.sum* ./

# Установка зависимостей
RUN go mod download

# Копирование исходного кода бэкенда
COPY backend/ .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -a -installsuffix cgo -o main ./cmd/main.go

# Этап 2: Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копирование собранного бэкенда
COPY --from=backend-builder /app/main .

# Копирование миграций
COPY db ./db

EXPOSE 8080

CMD ["/app/main"]

