package repository

import (
	"context"
	"payment-service/internal/model"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, p *model.Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*model.Payment, error) {
	var p model.Payment
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) GetByOrderID(ctx context.Context, orderID uint64) (*model.Payment, error) {
	var p model.Payment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) List(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
