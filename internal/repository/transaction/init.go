package transaction

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/MentalMentos/techFin/pkg/logger"
	"github.com/jackc/pgx/v4"
)

type Transaction interface {
	CreateTransaction(ctx context.Context, tx pgx.Tx, userID int, amount float64, targetUserID *int) error
	GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error)
}

type TransactionRepository struct {
	Transaction
}

func NewTransactionRepository(db db.Client, redis redis.IRedis, logger logger.Logger) *TransactionRepository {
	return &TransactionRepository{
		NewTransactionRepo(db, redis, logger),
	}
}
