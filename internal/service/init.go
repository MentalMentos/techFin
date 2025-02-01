package service

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/MentalMentos/techFin/internal/repository"
	"github.com/MentalMentos/techFin/pkg/logger"
)

type Service interface {
	GetBalance(ctx context.Context, userID int) (float64, error)
	UpdateBalance(ctx context.Context, userID int, amount float64) (float64, error)
	Transfer(ctx context.Context, fromUserID, toUserID int, amount float64) error
	GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error)
}

type ServiceImpl struct {
	repo      *repository.RepositoryImpl
	txManager db.TxManager
	logger    logger.Logger
}

func NewService(repo *repository.RepositoryImpl, tx db.TxManager, logger logger.Logger) *ServiceImpl {
	return &ServiceImpl{
		repo:      repo,
		txManager: tx,
		logger:    logger,
	}
}
