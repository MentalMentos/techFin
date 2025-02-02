package tests_test

import (
	"context"
	mock2 "github.com/MentalMentos/techFin/internal/mock"
	"github.com/MentalMentos/techFin/internal/repository/balance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateBalance(t *testing.T) {
	mockRedisClient := new(mock2.MockRedisClient)
	mockRedisClient.On("Set", mock.Anything, "balance:1", "150.0", mock.Anything).Return(nil)

	mockDbClient := new(mock2.MockPgClient) // Ваш мок для DbClient
	r := balance.NewBalanceRepo(mockDbClient, mockRedisClient, nil)

	mockTx := new(mock2.MockTx)
	mockTx.On("Begin", mock.Anything).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	mockTxManager := new(mock2.MockTxManager)
	mockTxManager.On("RepeatableRead", mock.Anything, mock.Anything).Return(nil)
	updatedBalance, err := r.UpdateBalance(context.Background(), mockTx, 1, 50.0)

	assert.NoError(t, err)
	mockTx.AssertExpectations(t)
	mockTxManager.AssertExpectations(t)
	assert.Equal(t, 150.0, updatedBalance)
	mockRedisClient.AssertExpectations(t)
}
