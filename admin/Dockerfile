FROM node:20-alpine AS builder

WORKDIR /app

# Аргументы для сборки
ARG VITE_API_URL=
ARG VITE_ADMIN_PASSWORD=admin
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_ADMIN_PASSWORD=$VITE_ADMIN_PASSWORD

# Копирование файлов зависимостей
COPY package*.json ./

# Установка зависимостей
RUN npm install

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN npm run build

FROM nginx:alpine

# Копирование собранного приложения
COPY --from=builder /app/dist /usr/share/nginx/html

# Копирование конфигурации nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

