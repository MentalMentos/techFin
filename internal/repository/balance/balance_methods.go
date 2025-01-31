package balance

import (
	"context"
	"fmt"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func (r *BalanceRepository) CreateBalance(ctx context.Context, tx pgx.Tx, userID int) error {
	_, err := tx.Exec(ctx, "INSERT INTO balances (user_id, balance VALUES ($1, $2) ON CONFLICT (user_id) DO NOTHING;", userID, 0)
	if err != nil {
		return errors.Wrap(err, "failed to create balance")
	}
	return nil
}

func (r *BalanceRepository) GetBalance(ctx context.Context, userID int) (float64, error) {
	balanceStr, err := r.redisClient.Get(ctx, fmt.Sprintf("balance:%d", userID))
	if err == nil {
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err == nil {
			return balance, nil
		}
	}

	var balance float64
	err = r.db.DB().QueryRowContext(ctx, db.Query{
		Name:     "get_balance",
		QueryRaw: "SELECT balance FROM users WHERE id=$1;",
	}, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	err = r.redisClient.Set(ctx, fmt.Sprintf("balance:%d", userID), fmt.Sprintf("%f", balance), 15*time.Minute)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *BalanceRepository) UpdateBalance(ctx context.Context, tx pgx.Tx, userID int, amount float64) (float64, error) {
	var updatedBalance float64

	_, err := tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2;", amount, userID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to update balance")
	}

	err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&updatedBalance)
	if err != nil {
		return 0, errors.Wrap(err, "failed to fetch updated balance")
	}

	err = r.redisClient.Set(ctx, fmt.Sprintf("balance:%d", userID), fmt.Sprintf("%f", updatedBalance), 15*time.Minute)
	if err != nil {
		return 0, errors.Wrap(err, "failed to cache balance in Redis")
	}

	return updatedBalance, nil
}
