package service

import (
	"context"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/MentalMentos/techFin/pkg/helpers"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
)

// GetBalance возвращает баланс пользователя, используя репозиторий
func (s *ServiceImpl) GetBalance(ctx context.Context, userID int) (float64, error) {
	// Получаем баланс через репозиторий
	balance, err := s.repo.BalanceRepository().GetBalance(ctx, userID)
	if err != nil {
		// Логируем ошибку при получении баланса
		s.logger.Info(helpers.ServicePrefix, helpers.ServiceGetBalanceError)
		return 0, errors.Wrap(err, "failed to get balance")
	}
	// Логируем успешное получение баланса
	s.logger.Info(helpers.ServicePrefix, "Balance retrieved successfully")
	return balance, nil
}

// UpdateBalance обновляет баланс пользователя с помощью транзакции
func (s *ServiceImpl) UpdateBalance(ctx context.Context, userID int, amount float64) (float64, error) {
	var updatedBalance float64

	// Операция с транзакцией для обновления баланса
	err := s.txManager.RepeatableRead(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		var err error
		// Обновляем баланс через репозиторий
		updatedBalance, err = s.repo.BalanceRepository().UpdateBalance(txCtx, tx, userID, amount)
		if err != nil {
			// Логируем ошибку при обновлении баланса
			s.logger.Info(helpers.ServicePrefix, helpers.ServiceUpdateBalanceError)
			return errors.Wrap(err, "failed to update balance")
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	// Логируем успешное обновление баланса
	s.logger.Info(helpers.ServicePrefix, "Balance updated successfully")
	return updatedBalance, nil
}

// Transfer выполняет перевод средств между двумя пользователями в рамках транзакции
func (s *ServiceImpl) Transfer(ctx context.Context, fromUserID, toUserID int, amount float64) error {
	// Получаем баланс отправителя
	balanceFromUser, err := s.repo.BalanceRepository().GetBalance(ctx, fromUserID)
	if err != nil {
		// Логируем ошибку при получении баланса отправителя
		s.logger.Info(helpers.ServicePrefix, helpers.ServiceGetBalanceError)
		return errors.Wrap(err, "failed to get balance")
	}
	log.Printf("Balance of user %d: %.2f, Transfer amount: %.2f", fromUserID, balanceFromUser, amount) // Логируем баланс и сумму перевода
	if amount <= 0 || balanceFromUser < amount {
		// Логируем ошибку, если сумма перевода некорректна
		s.logger.Info(helpers.ServicePrefix, helpers.ServiceInvalidAmount)
		return errors.New("amount must be positive")
	}

	// Осуществляем перевод в рамках транзакции
	return s.txManager.Serializable(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		// Обновляем баланс отправителя
		newBalanceFromUser, err := s.repo.BalanceRepository().UpdateBalance(txCtx, tx, fromUserID, -amount)
		if err != nil {
			// Логируем ошибку при обновлении баланса отправителя
			s.logger.Info(helpers.ServicePrefix, helpers.ServiceTransferError)
			return err
		}

		// Проверка на недостаточность средств
		if newBalanceFromUser < 0 {
			// Логируем ошибку, если недостаточно средств
			s.logger.Info(helpers.ServicePrefix, helpers.ServiceInsufficientFunds)
			return errors.New("insufficient funds")
		}

		// Обновляем баланс получателя
		if _, err := s.repo.BalanceRepository().UpdateBalance(txCtx, tx, toUserID, amount); err != nil {
			// Логируем ошибку при обновлении баланса получателя
			s.logger.Info(helpers.ServicePrefix, helpers.ServiceTransferError)
			return err
		}

		// Создаем транзакцию для отправителя
		if err := s.repo.TransactionRepository().CreateTransaction(txCtx, tx, fromUserID, -amount, &toUserID); err != nil {
			// Логируем ошибку при создании транзакции для отправителя
			s.logger.Info(helpers.ServicePrefix, helpers.ServiceTransactionError)
			return err
		}

		// Создаем транзакцию для получателя
		if err := s.repo.TransactionRepository().CreateTransaction(txCtx, tx, toUserID, amount, &fromUserID); err != nil {
			// Логируем ошибку при создании транзакции для получателя
			s.logger.Info(helpers.ServicePrefix, helpers.ServiceTransactionError)
			return err
		}
		// Логируем успешный перевод
		s.logger.Info(helpers.ServicePrefix, "Transfer completed successfully")
		return nil
	})
}

// GetLastTransactions получает последние транзакции пользователя
func (s *ServiceImpl) GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	// Осуществляем операцию с транзакцией для получения последних транзакций
	err := s.txManager.ReadCommitted(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		var err error
		// Получаем транзакции через репозиторий
		transactions, err = s.repo.TransactionRepository().GetLastTransactions(txCtx, userID)
		if err != nil {
			// Логируем ошибку при получении транзакций
			s.logger.Info(helpers.ServicePrefix, helpers.ServiceTransactionError)
			return errors.Wrap(err, "failed to get transactions")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Логируем успешное получение транзакций
	s.logger.Info(helpers.ServicePrefix, "Transactions retrieved successfully")
	return transactions, nil
}
