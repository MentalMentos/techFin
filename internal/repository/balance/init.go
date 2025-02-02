package balance

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/MentalMentos/techFin/pkg/logger"
	"github.com/jackc/pgx/v4"
)

type Balance interface {
	GetBalance(ctx context.Context, userID int) (float64, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, userID int, amount float64) (float64, error)
}

type BalanceRepository struct {
	Balance
}

func NewBalanceRepository(db db.Client, redis redis.IRedis, logger logger.Logger) BalanceRepository {
	return BalanceRepository{
		NewBalanceRepo(db, redis, logger),
	}
}
