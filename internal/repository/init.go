package repository

import (
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/MentalMentos/techFin/internal/repository/balance"
	"github.com/MentalMentos/techFin/internal/repository/transaction"
)

type Repository interface {
	BalanceRepository() balance.Balance
	TransactionRepository() transaction.Transaction
}

type RepositoryImpl struct {
	balanceRepo     balance.Balance
	transactionRepo transaction.Transaction
}

func NewRepository(dbClient db.Client, redisClient redis.IRedis) *RepositoryImpl {
	return &RepositoryImpl{
		balanceRepo:     balance.New(dbClient, redisClient),
		transactionRepo: transaction.NewTransactionRepository(dbClient, redisClient),
	}
}

func (r *RepositoryImpl) BalanceRepository() balance.Balance {
	return r.balanceRepo
}

func (r *RepositoryImpl) TransactionRepository() transaction.Transaction {
	return r.transactionRepo
}
