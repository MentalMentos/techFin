package repository

import (
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/MentalMentos/techFin/internal/repository/balance"
	"github.com/MentalMentos/techFin/internal/repository/transaction"
	"github.com/MentalMentos/techFin/pkg/logger"
)

type Repository interface {
	BalanceRepository() balance.Balance
	TransactionRepository() transaction.Transaction
}

type RepositoryImpl struct {
	balanceRepo     balance.Balance
	transactionRepo transaction.Transaction
	logger          logger.Logger
}

func NewRepository(dbClient db.Client, redisClient redis.IRedis, logger logger.Logger) *RepositoryImpl {
	return &RepositoryImpl{
		balanceRepo:     balance.New(dbClient, redisClient, logger),
		transactionRepo: transaction.NewTransactionRepository(dbClient, redisClient, logger),
	}
}

func (r *RepositoryImpl) BalanceRepository() balance.Balance {
	return r.balanceRepo
}

func (r *RepositoryImpl) TransactionRepository() transaction.Transaction {
	return r.transactionRepo
}
