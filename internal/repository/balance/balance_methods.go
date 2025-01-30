package balance

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/db/pg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (r *BalanceRepository) GetBalance(ctx *gin.Context, userID int) (float64, error) {
	var balance float64
	err := r.db.DB().QueryRowContext(ctx, db.Query{
		Name:     "get_balance",
		QueryRaw: "SELECT balance FROM users WHERE id=$1;",
	}, userID).Scan(&balance)
	return balance, err
}

func (r *BalanceRepository) UpdateBalance(ctx *gin.Context, tx db.TxManager, userID int, amount float64) error {
	return tx.ReadCommitted(ctx, func(txctx context.Context) error {
		tx, ok := txctx.Value(pg.TxKey).(pgx.Tx)
		if !ok {
			return errors.New("no transaction found in context")
		}
		_, err := tx.Exec(txctx, "UPDATE users SET balance = balance + $1 WHERE id = $2 RETURNING balance;", amount, userID)
		if err != nil {
			return errors.Wrap(err, "failed to update balance")
		}
		return nil
	})
}
