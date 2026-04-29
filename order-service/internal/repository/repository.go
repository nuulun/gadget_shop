package repository

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/model"
	"time"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("order not found")

type orderRow struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     uint64    `gorm:"column:user_id;not null;index"`
	Status     string    `gorm:"column:status;not null;default:'pending'"`
	TotalPrice float64   `gorm:"column:total_price;not null;default:0"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (orderRow) TableName() string { return "orders" }

type itemRow struct {
	ID        uint64  `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID   uint64  `gorm:"column:order_id;not null;index"`
	ProductID uint64  `gorm:"column:product_id;not null"`
	Quantity  int     `gorm:"column:quantity;not null;default:1"`
	Price     float64 `gorm:"column:price;not null"`
}

func (itemRow) TableName() string { return "order_items" }

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, o model.Order) (model.Order, error) {
	row := orderRow{UserID: o.UserID, Status: o.Status, TotalPrice: o.TotalPrice}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return model.Order{}, fmt.Errorf("create order: %w", err)
	}
	items := make([]itemRow, len(o.Items))
	for i, item := range o.Items {
		items[i] = itemRow{OrderID: row.ID, ProductID: item.ProductID, Quantity: item.Quantity, Price: item.Price}
	}
	if len(items) > 0 {
		r.db.WithContext(ctx).Create(&items)
	}
	o.ID = row.ID
	return r.GetByID(ctx, row.ID)
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (model.Order, error) {
	var row orderRow
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return model.Order{}, ErrNotFound }
		return model.Order{}, err
	}
	var items []itemRow
	r.db.WithContext(ctx).Where("order_id = ?", id).Find(&items)
	return rowToModel(row, items), nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uint64) ([]model.Order, error) {
	var rows []orderRow
	r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&rows)
	return r.enrichOrders(ctx, rows)
}

func (r *Repository) List(ctx context.Context) ([]model.Order, error) {
	var rows []orderRow
	r.db.WithContext(ctx).Order("created_at DESC").Find(&rows)
	return r.enrichOrders(ctx, rows)
}

func (r *Repository) enrichOrders(ctx context.Context, rows []orderRow) ([]model.Order, error) {
	orders := make([]model.Order, len(rows))
	for i, row := range rows {
		var items []itemRow
		r.db.WithContext(ctx).Where("order_id = ?", row.ID).Find(&items)
		orders[i] = rowToModel(row, items)
	}
	return orders, nil
}

func rowToModel(r orderRow, items []itemRow) model.Order {
	modelItems := make([]model.OrderItem, len(items))
	for i, it := range items {
		modelItems[i] = model.OrderItem{ID: it.ID, OrderID: it.OrderID, ProductID: it.ProductID, Quantity: it.Quantity, Price: it.Price}
	}
	return model.Order{ID: r.ID, UserID: r.UserID, Status: r.Status, TotalPrice: r.TotalPrice, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt, Items: modelItems}
}
