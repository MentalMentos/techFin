package response

import "github.com/MentalMentos/techFin/internal/models"

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

type TransactionsResponse struct {
	Transactions []models.Transaction `json:"transactions"`
}

type StandardResponse struct {
	Status  string      `json:"status"` // "success" | "error"
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
