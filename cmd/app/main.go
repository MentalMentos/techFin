package main

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db/pg"
	"github.com/MentalMentos/techFin/internal/clients/db/pg/transaction"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/MentalMentos/techFin/internal/clients/redis/go_redis"
	"github.com/MentalMentos/techFin/internal/handlers"
	"github.com/MentalMentos/techFin/internal/router"
	"github.com/MentalMentos/techFin/internal/service"
	"github.com/MentalMentos/techFin/pkg/helpers"
	zaplogger "github.com/MentalMentos/techFin/pkg/logger/zap"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
)

func main() {
	// Контекст приложения
	ctx := context.Background()

	// Инициализация логгера
	myLogger := zaplogger.New()

	// Загрузка переменных окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		// В случае ошибки загрузки конфигурации завершение работы с фатальным логом
		myLogger.Fatal(helpers.AppPrefix, "Failed to load env paths")
	}

	// Инициализация postgres клиента
	pgClient, err := pg.New(ctx)
	if err != nil {
		// В случае ошибки инициализации базы данных завершаем работу
		myLogger.Fatal(helpers.AppPrefix, "Failed to initialize database client")
	}
	defer pgClient.Close() // Закрытие соединения с базой данных по завершении работы

	// Инициализация менеджера транзакций
	txManager := transaction.NewTransactionManager(pgClient.DB())

	// Инициализация конфигурации Redis
	redisCfg, err := redis.NewRedisConfig()
	if err != nil {
		// В случае ошибки инициализации конфигурации Redis завершаем работу
		myLogger.Fatal(helpers.AppPrefix, "Failed to initialize Redis config")
	}

	// Инициализация клиента Redis
	redisClient, err := go_redis.NewGoRedisClient(redisCfg)
	if err != nil {
		// В случае ошибки инициализации Redis клиента завершаем работу
		myLogger.Fatal(helpers.AppPrefix, "Failed to initialize Redis client")
	}

	// Инициализация репозитория для работы с данными

	// Инициализация сервиса, который использует репозиторий и менеджер транзакций
	bankService := service.NewBankService(pgClient, redisClient, txManager, myLogger)

	// Инициализация логгера для обработчиков
	zap, err := zap.NewProduction()
	if err != nil {
		// В случае ошибки инициализации zap логгера завершаем работу
		myLogger.Fatal(helpers.AppPrefix, "Failed to initialize zap logger")
	}

	// Инициализация обработчиков HTTP-запросов
	bankHandlers := handlers.NewHandler(bankService, zap)

	// Настройка маршрутизатора для API
	r := router.SetupRouter(bankHandlers)

	// Запуск HTTP-сервера на порту 8080
	if err = r.Run("localhost:8080"); err != nil {
		// В случае ошибки запуска сервера выводим ошибку
		log.Fatalf("failed to start server: %v", err)
	}
}
