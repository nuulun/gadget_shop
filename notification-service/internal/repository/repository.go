package repository

import (
	"context"
	"notification-service/internal/model"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, n *model.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *Repository) List(ctx context.Context) ([]model.Notification, error) {
	var notifications []model.Notification
	if err := r.db.WithContext(ctx).Order("created_at desc").Limit(100).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
