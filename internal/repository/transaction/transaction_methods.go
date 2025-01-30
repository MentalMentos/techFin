package transaction

import (
	"context"
	"fmt"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"time"
)

func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx pgx.Tx, userID int, amount float64, transactionType string, targetUserID *int) error {
	_, err := tx.Exec(ctx, "INSERT INTO transactions (user_id, amount, transaction_type, target_user_id, status) VALUES ($1, $2, $3, $4, 'completed');",
		userID, amount, transactionType, targetUserID)
	if err != nil {
		tx.Rollback(ctx)
		return errors.Wrap(err, "failed to insert transaction")
	}

	transactionKey := fmt.Sprintf("transaction:%d:%d", userID, time.Now().Unix())
	transactionData := map[string]interface{}{
		"user_id":          userID,
		"amount":           amount,
		"transaction_type": transactionType,
		"target_user_id":   targetUserID,
		"status":           "completed",
	}

	err = r.redisClient.SetObject(ctx, transactionKey, transactionData, 24*time.Hour)
	if err != nil {
		return errors.Wrap(err, "failed to cache transaction in Redis")
	}

	return nil
}

func (r *TransactionRepository) GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	cacheKey := fmt.Sprintf("last_transactions:%d", userID)
	var cachedTransactions []models.Transaction
	err := r.redisClient.GetObject(ctx, cacheKey, &cachedTransactions)
	if err == nil && len(cachedTransactions) > 0 {
		return cachedTransactions, nil
	}

	rows, err := r.db.DB().QueryContext(ctx, db.Query{
		Name:     "get_last_transactions",
		QueryRaw: "SELECT id, amount, transaction_type, target_user_id, status, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT 10",
	}, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query transactions")
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.Type, &t.TargetUserID, &t.Status, &t.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan transaction")
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate over transaction rows")
	}

	err = r.redisClient.SetObject(ctx, cacheKey, transactions, 15*time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cache last transactions in Redis")
	}

	return transactions, nil
}
