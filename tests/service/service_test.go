// /test/service/service_test.go
package service_test

import (
	"context"
	"testing"

	"github.com/MentalMentos/techFin/internal/repository/mocks"
	"github.com/MentalMentos/techFin/internal/service"
	"github.com/MentalMentos/techFin/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetBalance(t *testing.T) {
	// Создаем моки для репозитория и логгера
	mockRepo := new(mocks.RepositoryImpl)
	mockBalanceRepo := new(mocks.BalanceRepository)
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLogger := new(mocks.Logger)

	// Настроим моки репозитория
	mockRepo.On("BalanceRepository").Return(mockBalanceRepo)
	mockRepo.On("TransactionRepository").Return(mockTransactionRepo)

	// Настроим мок для BalanceRepository
	mockBalanceRepo.On("GetBalance", mock.Anything, 1).Return(100.0, nil)

	// Создаем сервис с моками
	s := service.NewService(mockRepo, nil, mockLogger)

	// Вызовем метод
	balance, err := s.GetBalance(context.Background(), 1)

	// Проверим результаты
	assert.NoError(t, err)
	assert.Equal(t, 100.0, balance)

	// Проверим, что методы были вызваны
	mockRepo.AssertExpectations(t)
	mockBalanceRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestService_UpdateBalance(t *testing.T) {
	// Создаем моки для репозитория и логгера
	mockRepo := new(mocks.RepositoryImpl)
	mockBalanceRepo := new(mocks.BalanceRepository)
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLogger := new(mocks.Logger)

	// Настроим моки репозитория
	mockRepo.On("BalanceRepository").Return(mockBalanceRepo)
	mockRepo.On("TransactionRepository").Return(mockTransactionRepo)

	// Настроим мок для BalanceRepository
	mockBalanceRepo.On("UpdateBalance", mock.Anything, nil, 1, 50.0).Return(150.0, nil)

	// Создаем сервис с моками
	s := service.NewService(mockRepo, nil, mockLogger)

	// Вызовем метод
	newBalance, err := s.UpdateBalance(context.Background(), 1, 50.0)

	// Проверим результаты
	assert.NoError(t, err)
	assert.Equal(t, 150.0, newBalance)

	// Проверим, что методы были вызваны
	mockRepo.AssertExpectations(t)
	mockBalanceRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
