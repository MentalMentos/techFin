package balance

import (
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/gin-gonic/gin"
)

type Balance interface {
	GetBalance(ctx *gin.Context, userID int) (float64, error)
	UpdateBalance(ctx *gin.Context, tx db.TxManager, userID int, amount float64) (float64, error)
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
