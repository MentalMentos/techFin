include .env
LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

# Поднимаем контейнеры и запускаем приложение
up:
	docker-compose up --build -d
	go run cmd/app/main.go

# Останавливаем контейнеры
down:
	docker-compose down

# Устанавливаем зависимости
install-deps:
	GOBIN=$(CURDIR)/bin go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

# Полный запуск: поднимаем контейнеры (с миграциями) и сервер
run: up
