package balance

import (
	"awesomeProject3/internal/models"
	"github.com/gin-gonic/gin"
	//"github.com/jackc/pgx/v4/pgxpool"
)

func (repo *BalanceRepository) CreateBalance(ctx *gin.Context, user models.User) (int, error) {
	const mark = "[ REPOSITORY_CREATE_BALANCE ]"

	const query = "INSERT INTO users DEFAULT VALUES "

}

func (repo *BalanceRepository) GetBalance(ctx *gin.Context, user_id int) (models.User, error) {}
