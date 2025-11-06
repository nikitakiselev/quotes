# Инфраструктура

Эта директория содержит скрипты и конфигурацию для развертывания проекта.

## Деплой

Для ручного деплоя на сервер используйте скрипт:

```bash
bash infra/deploy.sh
```

Или выполните команды вручную:

```bash
cd ~/infra/services/quotes
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Автоматический деплой

Проект настроен на автоматический деплой через GitHub Actions при пуше в ветку `main`.

Требуется настроить секрет `SSH_PRIVATE_KEY` в настройках репозитория GitHub.

