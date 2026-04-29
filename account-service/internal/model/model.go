package model

import "time"

type User struct {
	ID         uint64    `json:"id"`
	Login      string    `json:"login"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
	Age        uint32    `json:"age"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	ID         uint64 `json:"id"`
	Login      string `json:"login"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Age        uint32 `json:"age"`
}

type UpdateUserInput struct {
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	MiddleName *string `json:"middle_name,omitempty"`
	Age        *uint32 `json:"age,omitempty"`
}
