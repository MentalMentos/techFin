package service

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/MentalMentos/techFin/internal/repository/balance"
	"github.com/MentalMentos/techFin/internal/repository/transaction"
	"github.com/MentalMentos/techFin/pkg/logger"
)

type Bank interface {
	GetBalance(ctx context.Context, userID int) (float64, error)
	UpdateBalance(ctx context.Context, userID int, amount float64) (float64, error)
	Transfer(ctx context.Context, fromUserID, toUserID int, amount float64) error
	GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error)
}

type BankService struct {
	Bank
}

func NewBankService(db db.Client, redis redis.IRedis, tx db.TxManager, logger logger.Logger) *BankService {
	return &BankService{
		NewService(transaction.NewTransactionRepository(db, redis, logger), balance.NewBalanceRepository(db, redis, logger), tx, logger),
	}
}
