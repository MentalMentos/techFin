package transaction

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/MentalMentos/techFin/internal/models"
)

type Transaction interface {
	CreateTransaction(ctx context.Context, tx db.TxManager, userID int, amount float64, transactionType string, targetUserID *int) error
	GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error)
}

type TransactionRepository struct {
	db          db.Client
	redisClient redis.IRedis
}

func NewTransactionRepository(db db.Client, redisClient redis.IRedis) *TransactionRepository {
	return &TransactionRepository{
		db:          db,
		redisClient: redisClient,
	}
}
