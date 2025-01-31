package request

//type CreateBalanceRequest struct {
//	UserID int `json:"user_id" binding:"required"`
//}

type UpdateBalanceRequest struct {
	UserID int     `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type TransferRequest struct {
	FromUserID int     `json:"from_user_id" binding:"required"`
	ToUserID   int     `json:"to_user_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
}

type GetTransactionsRequest struct {
	UserID int `json:"user_id" binding:"required"`
}
