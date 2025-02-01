package mocks

import (
	"github.com/MentalMentos/techFin/internal/repository/balance"
	"github.com/MentalMentos/techFin/internal/repository/transaction"
	"github.com/stretchr/testify/mock"
)

type RepositoryImpl struct {
	mock.Mock
}

func (r *RepositoryImpl) BalanceRepository() balance.Balance {
	args := r.Called()
	return args.Get(0).(balance.Balance)
}

func (r *RepositoryImpl) TransactionRepository() transaction.Transaction {
	args := r.Called()
	return args.Get(0).(transaction.Transaction)
}
