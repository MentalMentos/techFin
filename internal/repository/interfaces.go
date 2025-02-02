package repository

import (
	"context"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/jackc/pgx/v4"
)

//go:generate /Users/romchik/go/bin/mockgen -source=interfaces.go -destination=mocks/repository_mocks.go -package=mocks

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int) (float64, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, userID int, amount float64) (float64, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx pgx.Tx, userID int, amount float64, targetUserID *int) error
	GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error)
}

type TxManager interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
