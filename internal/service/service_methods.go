package service

import (
	"context"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (s *ServiceImpl) GetBalance(ctx context.Context, userID int) (float64, error) {
	balance, err := s.repo.BalanceRepository().GetBalance(ctx, userID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get balance")
	}
	return balance, nil
}

func (s *ServiceImpl) UpdateBalance(ctx context.Context, userID int, amount float64) (float64, error) {
	var updatedBalance float64

	err := s.txManager.RepeatableRead(ctx, func(txCtx context.Context) error {
		tx, ok := txCtx.Value("tx").(pgx.Tx)
		if !ok {
			return errors.New("failed to extract transaction from context")
		}

		var err error
		updatedBalance, err = s.repo.BalanceRepository().UpdateBalance(txCtx, tx, userID, amount)
		if err != nil {
			return errors.Wrap(err, "failed to update balance")
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return updatedBalance, nil
}

func (s *ServiceImpl) Transfer(ctx context.Context, fromUserID, toUserID int, amount float64) error {
	balanceFromUser, err := s.repo.BalanceRepository().GetBalance(ctx, fromUserID)
	if err != nil {
		return errors.Wrap(err, "failed to get balance")
	}
	if amount <= 0 || balanceFromUser < amount {
		return errors.New("amount must be positive")
	}

	return s.txManager.Serializable(ctx, func(txCtx context.Context) error {
		tx, ok := txCtx.Value("tx").(pgx.Tx)
		if !ok {
			return errors.New("failed to extract transaction from context")
		}
		newBalanceFromUser, err := s.repo.BalanceRepository().UpdateBalance(txCtx, tx, fromUserID, -amount)
		if err != nil {
			return err
		}

		if newBalanceFromUser < 0 {
			return errors.New("insufficient funds")
		}

		if _, err := s.repo.BalanceRepository().UpdateBalance(txCtx, tx, toUserID, amount); err != nil {
			return err
		}
		// Логируем операции перевода
		if err := s.repo.TransactionRepository().CreateTransaction(txCtx, tx, fromUserID, -amount, "transfer", &toUserID); err != nil {
			return err
		}
		if err := s.repo.TransactionRepository().CreateTransaction(txCtx, tx, toUserID, amount, "receive", &fromUserID); err != nil {
			return err
		}

		return nil
	})
}

func (s *ServiceImpl) GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		var err error
		// Передаем транзакцию в репозиторий для чтения
		transactions, err = s.repo.TransactionRepository().GetLastTransactions(txCtx, userID)
		if err != nil {
			return errors.Wrap(err, "failed to get transactions")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
