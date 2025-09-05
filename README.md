# Effective-Mobile
Тестовое задание Effective Mobile

Quick start:

Клонируйте репозиторий:
```bash
git clone https://github.com/Fedasov/Effective-Mobile.git
cd Effective-Mobile
```

Запустите сборку и запуск контейнеров:
```bash
docker compose up --build
```

Пример .env файла:
```bash
# Database configuration
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions

# Server configuration
SERVER_PORT=8080
```

Сервис будет доступен по адресу: http://localhost:8080
Swagger документация: http://localhost:8080/swagger/index.html