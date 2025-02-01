package models

import "time"

type Transaction struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Amount       float64   `json:"amount"`
	TargetUserID *int      `json:"target_user_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
