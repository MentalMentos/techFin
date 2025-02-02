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

	updatedBalance, err := r.UpdateBalance(context.Background(), nil, 1, 50.0)

	assert.NoError(t, err)
	assert.Equal(t, 150.0, updatedBalance)
	mockRedisClient.AssertExpectations(t)
}
