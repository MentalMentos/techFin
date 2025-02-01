package main

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db/pg"
	"github.com/MentalMentos/techFin/internal/clients/db/pg/transaction"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/MentalMentos/techFin/internal/clients/redis/go_redis"
	"github.com/MentalMentos/techFin/internal/handlers"
	"github.com/MentalMentos/techFin/internal/repository"
	"github.com/MentalMentos/techFin/internal/router"
	"github.com/MentalMentos/techFin/internal/service"
	"log"
)

func main() {
	ctx := context.Background()

	pgClient, err := pg.New(ctx)
	if err != nil {
		log.Fatalf("failed to create PostgreSQL client: %v", err)
	}
	defer pgClient.Close() // Закрытие соединения по завершении

	txManager := transaction.NewTransactionManager(pgClient.DB())

	redisCfg, err := redis.NewRedisConfig()
	redisClient := go_redis.NewGoRedisClient(redisCfg)
	if err != nil {
		log.Fatalf("failed to create Redis client: %v", err)
	}

	bankRepository := repository.NewRepository(pgClient, redisClient)
	bankService := service.NewService(bankRepository, txManager)
	bankHandlers := handlers.NewHandler(bankService)
	r := router.SetupRouter(bankHandlers)

	if err = r.Run("localhost:8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
