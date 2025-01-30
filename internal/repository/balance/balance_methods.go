package balance

import (
	"context"
	"fmt"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/db/pg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func (r *BalanceRepository) GetBalance(ctx *gin.Context, userID int) (float64, error) {
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

	err = r.redisClient.Set(ctx, fmt.Sprintf("balance:%d", userID), fmt.Sprintf("%f", balance), 10*time.Minute)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *BalanceRepository) UpdateBalance(ctx *gin.Context, tx db.TxManager, userID int, amount float64) (float64, error) {
	var updatedBalance float64
	err := tx.ReadCommitted(ctx, func(txctx context.Context) error {
		tx, ok := txctx.Value(pg.TxKey).(pgx.Tx)
		if !ok {
			return errors.New("no transaction found in context")
		}
		_, err := tx.Exec(txctx, "UPDATE users SET balance = balance + $1 WHERE id = $2 RETURNING balance;", amount, userID)
		if err != nil {
			return errors.Wrap(err, "failed to update balance")
		}

		err = tx.QueryRow(txctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&updatedBalance)
		if err != nil {
			return errors.Wrap(err, "failed to fetch updated balance")
		}

		err = r.redisClient.Set(ctx, fmt.Sprintf("balance:%d", userID), fmt.Sprintf("%f", updatedBalance), 10*time.Minute)
		if err != nil {
			return errors.Wrap(err, "failed to cache balance in Redis")
		}

		return nil
	})
	return updatedBalance, err
}
