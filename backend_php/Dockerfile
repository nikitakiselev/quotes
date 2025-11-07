# Build stage
FROM php:8.2-cli-alpine AS builder

# Install dependencies
RUN apk add --no-cache \
    git \
    unzip \
    postgresql-dev \
    linux-headers \
    $PHPIZE_DEPS

# Install PHP extensions including OPcache
RUN docker-php-ext-install pdo pdo_pgsql sockets opcache

# Install Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

# Copy application files
WORKDIR /app
COPY backend_php/composer.json backend_php/composer.lock* ./
RUN composer install --no-dev --optimize-autoloader --no-interaction

# Install RoadRunner binary
RUN composer require spiral/roadrunner-cli --no-interaction --prefer-dist && \
    php vendor/bin/rr get-binary -n && \
    mkdir -p /usr/local/bin && \
    cp /app/rr /usr/local/bin/rr && \
    chmod +x /usr/local/bin/rr

# Runtime stage
FROM php:8.2-cli-alpine

# Install runtime dependencies
RUN apk add --no-cache \
    postgresql-dev \
    linux-headers \
    $PHPIZE_DEPS \
    && docker-php-ext-install pdo pdo_pgsql sockets opcache \
    && apk del $PHPIZE_DEPS

# Configure OPcache for maximum performance
RUN echo "opcache.enable=1" >> /usr/local/etc/php/conf.d/opcache.ini && \
    echo "opcache.enable_cli=1" >> /usr/local/etc/php/conf.d/opcache.ini && \
    echo "opcache.memory_consumption=256" >> /usr/local/etc/php/conf.d/opcache.ini && \
    echo "opcache.interned_strings_buffer=16" >> /usr/local/etc/php/conf.d/opcache.ini && \
    echo "opcache.max_accelerated_files=10000" >> /usr/local/etc/php/conf.d/opcache.ini && \
    echo "opcache.validate_timestamps=0" >> /usr/local/etc/php/conf.d/opcache.ini && \
    echo "opcache.save_comments=0" >> /usr/local/etc/php/conf.d/opcache.ini && \
    echo "opcache.fast_shutdown=1" >> /usr/local/etc/php/conf.d/opcache.ini

# Copy application from builder
COPY --from=builder /app/vendor /app/vendor
COPY --from=builder /usr/local/bin/rr /usr/local/bin/rr
COPY backend_php /app

WORKDIR /app

# Expose port
EXPOSE 8082

# Run RoadRunner
CMD ["/usr/local/bin/rr", "serve", "-c", ".rr.yaml"]

