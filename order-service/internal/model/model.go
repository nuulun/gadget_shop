package model

import "time"

type Order struct {
	ID         uint64      `json:"id"`
	UserID     uint64      `json:"user_id"`
	Status     string      `json:"status"`
	TotalPrice float64     `json:"total_price"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Items      []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID        uint64  `json:"id"`
	OrderID   uint64  `json:"order_id"`
	ProductID uint64  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CreateOrderInput struct {
	UserID uint64            `json:"user_id"`
	Items  []CreateItemInput `json:"items"`
}

type CreateItemInput struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// ProductInfo is fetched from product-service.
type ProductInfo struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
