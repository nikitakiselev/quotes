# Этап 1: Сборка фронтенда
FROM node:20-alpine AS builder

WORKDIR /app

# Копирование файлов зависимостей
COPY package*.json ./

# Установка зависимостей
# Игнорируем postinstall скрипты, которые могут требовать patch-package
RUN npm install --ignore-scripts

# Копирование исходного кода
COPY . .

# Сборка фронтенда
# VITE_API_URL - URL backend для API запросов
# Если пусто, используется относительный путь (тот же домен)
ARG VITE_API_URL=
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

# Этап 2: Caddy для отдачи статики
FROM caddy:2-alpine

# Копирование собранного фронтенда
COPY --from=builder /app/dist /usr/share/caddy

# Копирование конфигурации Caddy
COPY Caddyfile /etc/caddy/Caddyfile

EXPOSE 80

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile"]
