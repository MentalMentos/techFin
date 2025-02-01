package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type TransactionRepository struct {
	mock.Mock
}

func (m *TransactionRepository) CreateTransaction(ctx context.Context, tx interface{}, userID int, amount float64, targetUserID *int) error {
	args := m.Called(ctx, tx, userID, amount, targetUserID)
	return args.Error(0)
}

func (m *TransactionRepository) GetLastTransactions(ctx context.Context, userID int) ([]interface{}, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]interface{}), args.Error(1)
}
