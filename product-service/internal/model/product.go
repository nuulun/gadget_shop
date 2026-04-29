package model

import "time"

// Product is the domain model – independent of DB or HTTP representation.
type Product struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	Image       string    `json:"image"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductInput is the DTO for product creation.
type CreateProductInput struct {
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
}

// UpdateProductInput is the DTO for product updates (all fields optional).
type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Brand       *string  `json:"brand,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
	Category    *string  `json:"category,omitempty"`
	Image       *string  `json:"image,omitempty"`
}

// ListFilter holds optional query parameters for listing products.
type ListFilter struct {
	Category string
	MinPrice float64
	MaxPrice float64
	Limit    int
	Offset   int
}
