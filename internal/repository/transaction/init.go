package transaction

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/redis/go-redis/v9"
)

type Transaction interface {
	CreateTransaction(ctx context.Context, tx db.TxManager, userID int, amount float64, transactionType string, targetUserID *int) error
	GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error)
}

type TransactionRepository struct {
	db          db.Client
	redisClient redis.Client
}
