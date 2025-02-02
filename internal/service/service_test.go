package service_test

import (
	"context"
	"github.com/MentalMentos/techFin/internal/service"
	"github.com/MentalMentos/techFin/pkg/logger"
	repo_mocks "github.com/MentalMentos/techFin/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceImpl_Transfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока для репозиториев
	balRepo := repo_mocks.NewBalanceRepository(ctrl)
	transRepo := repo_mocks.NewTransactionRepository(ctrl)
	//txManager := mocks.NewMockTxManager(ctrl)
	s := service.NewService(transRepo, balRepo, txManager, logger.New())

	ctx := context.Background()
	fromUserID, toUserID := 1, 2
	amount := 50.0

	t.Run("Success", func(t *testing.T) {
		balRepo.EXPECT().GetBalance(ctx, fromUserID).Return(100.0, nil)
		transRepo.EXPECT().CreateTransaction(ctx, gomock.Any(), fromUserID, -amount, gomock.Any()).Return(nil)
		transRepo.EXPECT().CreateTransaction(ctx, gomock.Any(), toUserID, amount, gomock.Any()).Return(nil)

		err := s.Transfer(ctx, fromUserID, toUserID, amount)
		assert.NoError(t, err)
	})
}
