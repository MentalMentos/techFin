package balance

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/jackc/pgx/v4"
)

type Balance interface {
	GetBalance(ctx context.Context, userID int) (float64, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, userID int, amount float64) (float64, error)
}

type BalanceRepository struct {
	db          db.Client
	redisClient redis.IRedis
}

func New(db db.Client, redisClient redis.IRedis) *BalanceRepository {
	return &BalanceRepository{
		db:          db,
		redisClient: redisClient,
	}
}
