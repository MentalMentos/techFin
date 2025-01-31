package service

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/MentalMentos/techFin/internal/repository"
)

type Service interface {
	CreateBalance(ctx context.Context, userID int) error
	GetBalance(ctx context.Context, userID int) (float64, error)
	UpdateBalance(ctx context.Context, userID int, amount float64) (float64, error)
	Transfer(ctx context.Context, fromUserID, toUserID int, amount float64) error
	GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error)
}

type ServiceImpl struct {
	repo      *repository.RepositoryImpl
	txManager db.TxManager
}

func NewService(repo *repository.RepositoryImpl, tx db.TxManager) *ServiceImpl {
	return &ServiceImpl{
		repo:      repo,
		txManager: tx,
	}
}
