package handlers

import (
	"github.com/MentalMentos/techFin/internal/data/request"
	"github.com/MentalMentos/techFin/internal/data/response"
	"github.com/MentalMentos/techFin/pkg/helpers"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/MentalMentos/techFin/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.BankService // Сервис для обработки логики
	logger  *zap.Logger          // Логгер для логирования запросов и ошибок
}

// NewHandler создает новый обработчик API
func NewHandler(service *service.BankService, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GetBalanceHandler - обработчик для получения баланса пользователя
func (h *Handler) GetBalanceHandler(c *gin.Context) {
	// Получаем user_id из URL параметра
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		// Логируем ошибку при некорректном user_id
		h.logger.Warn(helpers.HandlerPrefix+"Invalid user_id", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.StandardResponse{
			Status:  "error",
			Message: "invalid user_id",
		})
		return
	}

	// Получаем баланс через сервис
	balance, err := h.service.GetBalance(c.Request.Context(), userID)
	if err != nil {
		// Логируем ошибку при получении баланса
		h.logger.Error(helpers.HandlerPrefix+"Failed to get balance", zap.Int("user_id", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.StandardResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Логируем успешное получение баланса
	h.logger.Info(helpers.HandlerPrefix+"Balance retrieved", zap.Int("user_id", userID), zap.Float64("balance", balance))
	c.JSON(http.StatusOK, response.StandardResponse{
		Status:  "success",
		Message: "balance retrieved",
		Data:    response.BalanceResponse{Balance: balance},
	})
}

// UpdateBalanceHandler - обработчик для обновления баланса пользователя
func (h *Handler) UpdateBalanceHandler(c *gin.Context) {
	// Приводим тело запроса к соответствующей структуре
	var req request.UpdateBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Логируем ошибку при неверном теле запроса
		h.logger.Warn(helpers.HandlerPrefix+"Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.StandardResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}

	// Обновляем баланс через сервис
	updatedBalance, err := h.service.UpdateBalance(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		// Логируем ошибку при обновлении баланса
		h.logger.Error(helpers.HandlerPrefix+"Failed to update balance", zap.Int("user_id", req.UserID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.StandardResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Логируем успешное обновление баланса
	h.logger.Info(helpers.HandlerPrefix+"Balance updated", zap.Int("user_id", req.UserID), zap.Float64("new_balance", updatedBalance))
	c.JSON(http.StatusOK, response.StandardResponse{
		Status:  "success",
		Message: "balance updated",
		Data:    response.BalanceResponse{Balance: updatedBalance},
	})
}

// TransferHandler - обработчик для перевода денег между пользователями
func (h *Handler) TransferHandler(c *gin.Context) {
	// Приводим тело запроса к соответствующей структуре
	var req request.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Логируем ошибку при неверном теле запроса
		h.logger.Warn(helpers.HandlerPrefix+"Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.StandardResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}

	// Вызываем сервис для выполнения перевода
	if err := h.service.Transfer(c.Request.Context(), req.FromUserID, req.ToUserID, req.Amount); err != nil {
		// Логируем ошибку при ошибке перевода
		h.logger.Error(helpers.HandlerPrefix+"Failed to transfer funds", zap.Int("from_user_id", req.FromUserID), zap.Int("to_user_id", req.ToUserID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.StandardResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Логируем успешный перевод
	h.logger.Info(helpers.HandlerPrefix+"Transfer successful", zap.Int("from_user_id", req.FromUserID), zap.Int("to_user_id", req.ToUserID), zap.Float64("amount", req.Amount))
	c.JSON(http.StatusOK, response.StandardResponse{
		Status:  "success",
		Message: "transfer successful",
	})
}

// GetLastTransactionsHandler - обработчик для получения последних 10 транзакций пользователя
func (h *Handler) GetLastTransactionsHandler(c *gin.Context) {
	var req request.GetTransactionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Логируем ошибку при некорректном user_id
		h.logger.Warn(helpers.HandlerPrefix+"Invalid user_id", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.StandardResponse{
			Status:  "error",
			Message: "invalid user_id",
		})
		return
	}

	// Получаем последние транзакции через сервис
	transactions, err := h.service.GetLastTransactions(c.Request.Context(), req.UserID)
	if err != nil {
		// Логируем ошибку при получении транзакций
		h.logger.Error(helpers.HandlerPrefix+"Failed to get last transactions", zap.Int("user_id", req.UserID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.StandardResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Логируем успешное получение транзакций
	h.logger.Info(helpers.HandlerPrefix+"Transactions retrieved", zap.Int("user_id", req.UserID), zap.Int("count", len(transactions)))
	c.JSON(http.StatusOK, response.StandardResponse{
		Status:  "success",
		Message: "transactions retrieved",
		Data:    response.TransactionsResponse{Transactions: transactions},
	})
}
