package handlers

import (
	"github.com/MentalMentos/techFin/internal/data/request"
	"github.com/MentalMentos/techFin/internal/data/response"
	"net/http"
	"strconv"

	"github.com/MentalMentos/techFin/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.ServiceImpl
}

// NewHandler создает новый обработчик API
func NewHandler(service *service.ServiceImpl) *Handler {
	return &Handler{service: service}
}

// GetBalanceHandler - получение баланса пользователя
func (h *Handler) GetBalanceHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.StandardResponse{
			Status:  "error",
			Message: "invalid user_id",
		})
		return
	}

	balance, err := h.service.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.StandardResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.StandardResponse{
		Status:  "success",
		Message: "balance retrieved",
		Data:    response.BalanceResponse{Balance: balance},
	})
}

// UpdateBalanceHandler - обновление баланса пользователя
func (h *Handler) UpdateBalanceHandler(c *gin.Context) {
	var req request.UpdateBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.StandardResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}

	updatedBalance, err := h.service.UpdateBalance(c, req.UserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.StandardResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.StandardResponse{
		Status:  "success",
		Message: "balance updated",
		Data:    response.BalanceResponse{Balance: updatedBalance},
	})
}

// TransferHandler - перевод денег между пользователями
func (h *Handler) TransferHandler(c *gin.Context) {
	var req request.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.StandardResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}

	if err := h.service.Transfer(c, req.FromUserID, req.ToUserID, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, response.StandardResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.StandardResponse{
		Status:  "success",
		Message: "transfer successful",
	})
}

// GetLastTransactionsHandler - получение последних 10 транзакций пользователя
func (h *Handler) GetLastTransactionsHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.StandardResponse{
			Status:  "error",
			Message: "invalid user_id",
		})
		return
	}

	transactions, err := h.service.GetLastTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.StandardResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.StandardResponse{
		Status:  "success",
		Message: "transactions retrieved",
		Data:    response.TransactionsResponse{Transactions: transactions},
	})
}
