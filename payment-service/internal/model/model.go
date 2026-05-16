package model

import "time"

type Payment struct {
	ID        uint64    `json:"id"`
	OrderID   uint64    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePaymentInput struct {
	OrderID uint64  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
}
