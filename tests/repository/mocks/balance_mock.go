package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type BalanceRepository struct {
	mock.Mock
}

func (m *BalanceRepository) GetBalance(ctx context.Context, userID int) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *BalanceRepository) UpdateBalance(ctx context.Context, tx interface{}, userID int, amount float64) (float64, error) {
	args := m.Called(ctx, tx, userID, amount)
	return args.Get(0).(float64), args.Error(1)
}
