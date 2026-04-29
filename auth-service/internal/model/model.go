package model

import "time"

type User struct {
	ID           uint64
	Login        string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	ID        uint64
	UserID    uint64
	Token     string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type RegisterInput struct {
	ID       uint64 `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	LoginOrEmail string `json:"login_or_email"`
	Password     string `json:"password"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
