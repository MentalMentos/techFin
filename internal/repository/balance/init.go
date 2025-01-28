package balance

import (
	"awesomeProject3/internal/clients/db"
	"awesomeProject3/internal/clients/redis"
	"awesomeProject3/internal/models"
	"github.com/gin-gonic/gin"
)

type Balance interface {
	CreateBalance(ctx *gin.Context, user models.User) error
	GetBalance(ctx *gin.Context, user_id int) (models.User, error)
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
