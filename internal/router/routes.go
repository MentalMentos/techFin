package router

import (
	"github.com/MentalMentos/techFin/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *handlers.Handler) *gin.Engine {
	r := gin.Default()

	// Группа API v1
	api := r.Group("/api")
	{
		//api.GET("/balance", handler.GetBalanceHandler)            // Получение баланса
		api.POST("/balance-update", handler.UpdateBalanceHandler) // Пополнение баланса

		api.POST("/transfer", handler.TransferHandler)               // Перевод денег
		api.GET("/transactions", handler.GetLastTransactionsHandler) // Последние 10 транзакций
	}

	return r
}
