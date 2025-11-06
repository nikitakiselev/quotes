# Этап 1: Сборка фронтенда
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Копирование файлов зависимостей фронтенда
COPY frontend/package*.json ./

# Установка зависимостей
RUN npm install

# Копирование исходного кода фронтенда
COPY frontend/ .

# Сборка фронтенда
# VITE_API_URL оставляем пустым для использования относительных путей
ARG VITE_API_URL=
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

# Этап 2: Сборка бэкенда
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

# Этап 3: Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копирование собранного бэкенда
COPY --from=backend-builder /app/main .

# Копирование собранного фронтенда
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Копирование миграций
COPY db ./db

EXPOSE 8080

ENV STATIC_DIR=/app/frontend/dist

CMD ["/app/main"]

