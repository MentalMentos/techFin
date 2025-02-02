package service_test

import (
	"context"
	"errors"
	"github.com/MentalMentos/techFin/internal/repository"
	"testing"

	"github.com/MentalMentos/techFin/internal/repository/mocks"
	"github.com/MentalMentos/techFin/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServiceImpl_Transfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание моков без указателей
	balRepo := mocks.NewMockBalanceRepository(ctrl)
	transRepo := mocks.NewMockTransactionRepository(ctrl)
	txManager := mocks.NewMockTxManager(ctrl)

	s := service.NewService(repository.TransactionRepository(),
		repository.BalanceRepository(balRepo),
		repository.TxManager(txManager),
		nil)

	ctx := context.Background()
	fromUserID, toUserID := 1, 2
	amount := 50.0

	t.Run("Успешный перевод", func(t *testing.T) {
		// Настройка моков
		txManager.EXPECT().Begin(ctx).Return(ctx, nil)
		balRepo.EXPECT().GetBalance(ctx, fromUserID).Return(100.0, nil)
		transRepo.EXPECT().CreateTransaction(ctx, gomock.Any(), fromUserID, -amount, gomock.Any()).Return(nil)
		transRepo.EXPECT().CreateTransaction(ctx, gomock.Any(), toUserID, amount, gomock.Any()).Return(nil)
		txManager.EXPECT().Commit(ctx).Return(nil)

		err := s.Transfer(ctx, fromUserID, toUserID, amount)
		assert.NoError(t, err)
	})

	t.Run("Недостаточно средств", func(t *testing.T) {
		txManager.EXPECT().Begin(ctx).Return(ctx, nil)
		balRepo.EXPECT().GetBalance(ctx, fromUserID).Return(30.0, nil)
		txManager.EXPECT().Rollback(ctx).Return(nil)

		err := s.Transfer(ctx, fromUserID, toUserID, amount)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient funds")
	})

	t.Run("Ошибка получения баланса", func(t *testing.T) {
		txManager.EXPECT().Begin(ctx).Return(ctx, nil)
		balRepo.EXPECT().GetBalance(ctx, fromUserID).Return(0.0, errors.New("database error"))
		txManager.EXPECT().Rollback(ctx).Return(nil)

		err := s.Transfer(ctx, fromUserID, toUserID, amount)
		assert.Error(t, err)
	})

	t.Run("Ошибка создания транзакции списания", func(t *testing.T) {
		txManager.EXPECT().Begin(ctx).Return(ctx, nil)
		balRepo.EXPECT().GetBalance(ctx, fromUserID).Return(100.0, nil)
		transRepo.EXPECT().CreateTransaction(ctx, gomock.Any(), fromUserID, -amount, gomock.Any()).
			Return(errors.New("transaction error"))
		txManager.EXPECT().Rollback(ctx).Return(nil)

		err := s.Transfer(ctx, fromUserID, toUserID, amount)
		assert.Error(t, err)
	})

	t.Run("Ошибка создания транзакции зачисления", func(t *testing.T) {
		txManager.EXPECT().Begin(ctx).Return(ctx, nil)
		balRepo.EXPECT().GetBalance(ctx, fromUserID).Return(100.0, nil)
		transRepo.EXPECT().CreateTransaction(ctx, gomock.Any(), fromUserID, -amount, gomock.Any()).Return(nil)
		transRepo.EXPECT().CreateTransaction(ctx, gomock.Any(), toUserID, amount, gomock.Any()).
			Return(errors.New("transaction error"))
		txManager.EXPECT().Rollback(ctx).Return(nil)

		err := s.Transfer(ctx, fromUserID, toUserID, amount)
		assert.Error(t, err)
	})

	t.Run("Перевод самому себе", func(t *testing.T) {
		err := s.Transfer(ctx, fromUserID, fromUserID, amount)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot transfer to the same account")
	})

	t.Run("Отрицательная сумма", func(t *testing.T) {
		err := s.Transfer(ctx, fromUserID, toUserID, -amount)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "amount must be positive")
	})
}
