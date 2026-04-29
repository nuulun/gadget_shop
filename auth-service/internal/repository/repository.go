package repository

import (
	"auth-service/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type userRow struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Login        string    `gorm:"column:login;uniqueIndex;not null"`
	Email        string    `gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (userRow) TableName() string { return "users" }

type refreshTokenRow struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint64     `gorm:"column:user_id;uniqueIndex;not null"`
	Token     string     `gorm:"column:token;uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (refreshTokenRow) TableName() string { return "refresh_tokens" }

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreateUser(ctx context.Context, u model.User) error {
	row := userRow{ID: u.ID, Login: u.Login, Email: u.Email, PasswordHash: u.PasswordHash}
	return r.db.WithContext(ctx).Create(&row).Error
}

func (r *Repository) GetUserByLoginOrEmail(ctx context.Context, s string) (model.User, error) {
	var row userRow
	err := r.db.WithContext(ctx).Where("login = ? OR email = ?", s, s).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, ErrNotFound
	}
	return userRowToModel(row), err
}

func (r *Repository) GetUserByID(ctx context.Context, id uint64) (model.User, error) {
	var row userRow
	err := r.db.WithContext(ctx).First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, ErrNotFound
	}
	return userRowToModel(row), err
}

func (r *Repository) SaveRefreshToken(ctx context.Context, rt model.RefreshToken) error {
	// Upsert: delete old token for user then insert new
	r.db.WithContext(ctx).Where("user_id = ?", rt.UserID).Delete(&refreshTokenRow{})
	row := refreshTokenRow{UserID: rt.UserID, Token: rt.Token, ExpiresAt: rt.ExpiresAt}
	return r.db.WithContext(ctx).Create(&row).Error
}

func (r *Repository) GetRefreshToken(ctx context.Context, token string) (model.RefreshToken, error) {
	var row refreshTokenRow
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.RefreshToken{}, ErrNotFound
	}
	return refreshRowToModel(row), err
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&refreshTokenRow{}).
		Where("token = ?", token).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}

func (r *Repository) DeleteUser(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Where("user_id = ?", id).Delete(&refreshTokenRow{})
		res := tx.Where("id = ?", id).Delete(&userRow{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("user not found")
		}
		return nil
	})
}

func userRowToModel(r userRow) model.User {
	return model.User{ID: r.ID, Login: r.Login, Email: r.Email, PasswordHash: r.PasswordHash, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
}

func refreshRowToModel(r refreshTokenRow) model.RefreshToken {
	return model.RefreshToken{ID: r.ID, UserID: r.UserID, Token: r.Token, ExpiresAt: r.ExpiresAt, RevokedAt: r.RevokedAt, CreatedAt: r.CreatedAt}
}
