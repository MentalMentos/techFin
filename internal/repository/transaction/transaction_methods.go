package transaction

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/db/pg"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx db.TxManager, userID int, amount float64, transactionType string, targetUserID *int) error {
	return tx.ReadCommitted(ctx, func(txctx context.Context) error {
		tx, ok := txctx.Value(pg.TxKey).(pgx.Tx)
		if !ok {
			return errors.New("no transaction found in context")
		}
		_, err := tx.Exec(txctx, "INSERT INTO transactions (user_id, amount, transaction_type, target_user_id, status) VALUES ($1, $2, $3, $4, 'completed');",
			userID, amount, transactionType, targetUserID)
		if err != nil {
			return errors.Wrap(err, "failed to update balance")
		}
		return nil
	})
}

func (r *TransactionRepository) GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	// Выполняем запрос для получения последних 10 транзакций
	rows, err := r.db.DB().QueryContext(ctx, db.Query{
		Name:     "get_last_transactions",
		QueryRaw: "SELECT id, amount, transaction_type, target_user_id, status, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT 10",
	}, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query transactions")
	}
	defer rows.Close()

	// Чтение данных из строк результата
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.Type, &t.TargetUserID, &t.Status, &t.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan transaction")
		}
		transactions = append(transactions, t)
	}

	// Проверка на ошибку после завершения чтения строк
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate over transaction rows")
	}

	return transactions, nil
}
