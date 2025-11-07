# Build stage
FROM php:8.2-cli-alpine AS builder

# Install dependencies
RUN apk add --no-cache \
    git \
    unzip \
    postgresql-dev \
    $PHPIZE_DEPS

# Install PHP extensions
RUN docker-php-ext-install pdo pdo_pgsql sockets

# Install Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

# Copy application files
WORKDIR /app
COPY backend_php/composer.json backend_php/composer.lock* ./
RUN composer install --no-dev --optimize-autoloader --no-interaction

# Install RoadRunner binary
RUN composer require spiral/roadrunner-cli --no-interaction --prefer-dist

# Runtime stage
FROM php:8.2-cli-alpine

# Install runtime dependencies
RUN apk add --no-cache \
    postgresql-dev \
    $PHPIZE_DEPS \
    && docker-php-ext-install pdo pdo_pgsql sockets \
    && apk del $PHPIZE_DEPS

# Copy application from builder
COPY --from=builder /app/vendor /app/vendor
COPY backend_php /app

WORKDIR /app

# Expose port
EXPOSE 8082

# Run RoadRunner
CMD ["php", "vendor/bin/rr", "serve", "-c", ".rr.yaml"]

