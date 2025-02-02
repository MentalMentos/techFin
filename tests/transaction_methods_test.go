package tests_test

import (
	"context"
	mock2 "github.com/MentalMentos/techFin/internal/mock"
	"testing"

	"github.com/MentalMentos/techFin/internal/repository/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction(t *testing.T) {
	mockRedisClient := new(mock2.MockRedisClient)
	mockRedisClient.On("SetObject", mock.Anything, "transaction:1:12345", mock.Anything, mock.Anything).Return(nil)

	mockDbClient := new(mock2.MockPgClient) // Ваш мок для DbClient
	tr := transaction.NewTransactionRepo(mockDbClient, mockRedisClient, nil)

	err := tr.CreateTransaction(context.Background(), nil, 1, 50.0, nil)

	assert.NoError(t, err)
	mockRedisClient.AssertExpectations(t)
}

func TestGetLastTransactions(t *testing.T) {
	mockRedisClient := new(mock2.MockRedisClient)
	mockRedisClient.On("GetObject", mock.Anything, mock.Anything).Return(nil)

	mockDbClient := new(mock2.MockPgClient) // Ваш мок для DbClient
	tr := transaction.NewTransactionRepo(mockDbClient, mockRedisClient, nil)

	transactions, err := tr.GetLastTransactions(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	mockRedisClient.AssertExpectations(t)
}
