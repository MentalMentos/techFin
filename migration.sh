#!/bin/bash
set -e  # Скрипт завершится при первой ошибке
source .env

export MIGRATION_DSN="host=pg port=5432 dbname=$PG_DATABASE_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable"

echo "Starting migrations..."
echo "Using DSN: $MIGRATION_DSN"
echo "Migration directory: $MIGRATION_DIR"

sleep 10  # Подождать, пока PostgreSQL будет доступен
goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v