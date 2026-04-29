package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"product-service/internal/model"

	"gorm.io/gorm"
)

// ErrNotFound is returned when a product cannot be found.
var ErrNotFound = errors.New("product not found")

// productRow is the GORM DB model – kept internal to this package.
type productRow struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;not null"`
	Brand       string    `gorm:"column:brand"`
	Description string    `gorm:"column:description"`
	Price       float64   `gorm:"column:price;not null"`
	Stock       int       `gorm:"column:stock;default:0"`
	Category    string    `gorm:"column:category"`
	Image       string    `gorm:"column:image_url"`
	IsActive    bool      `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (productRow) TableName() string { return "products" }

// Repository provides CRUD access to the products table.
type Repository struct {
	db *gorm.DB
}

// New creates a Repository.
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ─── Queries ──────────────────────────────────────────────────────────────────

// List returns a filtered, paginated slice of products.
func (r *Repository) List(ctx context.Context, f model.ListFilter) ([]model.Product, error) {
	q := r.db.WithContext(ctx).Model(&productRow{}).Where("is_active = true")

	if f.Category != "" {
		q = q.Where("category = ?", f.Category)
	}
	if f.MinPrice > 0 {
		q = q.Where("price >= ?", f.MinPrice)
	}
	if f.MaxPrice > 0 {
		q = q.Where("price <= ?", f.MaxPrice)
	}

	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var rows []productRow
	if err := q.Limit(limit).Offset(f.Offset).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("repository.List: %w", err)
	}

	products := make([]model.Product, len(rows))
	for i, row := range rows {
		products[i] = rowToModel(row)
	}
	return products, nil
}

// Count returns the total number of active products matching the filter.
func (r *Repository) Count(ctx context.Context, f model.ListFilter) (int64, error) {
	q := r.db.WithContext(ctx).Model(&productRow{}).Where("is_active = true")
	if f.Category != "" {
		q = q.Where("category = ?", f.Category)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("repository.Count: %w", err)
	}
	return count, nil
}

// GetByID returns a single product or ErrNotFound.
func (r *Repository) GetByID(ctx context.Context, id uint64) (model.Product, error) {
	var row productRow
	err := r.db.WithContext(ctx).First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Product{}, ErrNotFound
	}
	if err != nil {
		return model.Product{}, fmt.Errorf("repository.GetByID: %w", err)
	}
	return rowToModel(row), nil
}

// ─── Mutations ────────────────────────────────────────────────────────────────

// Create inserts a new product and returns it with its DB-assigned ID.
func (r *Repository) Create(ctx context.Context, in model.CreateProductInput) (model.Product, error) {
	row := productRow{
		Name:        in.Name,
		Brand:       in.Brand,
		Description: in.Description,
		Price:       in.Price,
		Stock:       in.Stock,
		Category:    in.Category,
		Image:       in.Image,
		IsActive:    true,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return model.Product{}, fmt.Errorf("repository.Create: %w", err)
	}
	return rowToModel(row), nil
}

// Update applies partial updates to a product.
func (r *Repository) Update(ctx context.Context, id uint64, in model.UpdateProductInput) (model.Product, error) {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Brand != nil {
		updates["brand"] = *in.Brand
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}
	if in.Price != nil {
		updates["price"] = *in.Price
	}
	if in.Stock != nil {
		updates["stock"] = *in.Stock
	}
	if in.Category != nil {
		updates["category"] = *in.Category
	}
	if in.Image != nil {
		updates["image_url"] = *in.Image
	}

	result := r.db.WithContext(ctx).Model(&productRow{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return model.Product{}, fmt.Errorf("repository.Update: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.Product{}, ErrNotFound
	}
	return r.GetByID(ctx, id)
}

// Delete soft-deletes a product by marking it inactive.
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Model(&productRow{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()})
	if result.Error != nil {
		return fmt.Errorf("repository.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// IsEmpty reports whether the products table has zero active rows.
func (r *Repository) IsEmpty(ctx context.Context) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&productRow{}).Where("is_active = true").Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

// BulkCreate inserts many products efficiently.
func (r *Repository) BulkCreate(ctx context.Context, inputs []model.CreateProductInput) error {
	rows := make([]productRow, len(inputs))
	for i, in := range inputs {
		rows[i] = productRow{
			Name:        in.Name,
			Brand:       in.Brand,
			Description: in.Description,
			Price:       in.Price,
			Stock:       in.Stock,
			Category:    in.Category,
			Image:       in.Image,
			IsActive:    true,
		}
	}
	return r.db.WithContext(ctx).CreateInBatches(rows, 100).Error
}

// ─── Mapper ───────────────────────────────────────────────────────────────────

func rowToModel(r productRow) model.Product {
	return model.Product{
		ID:          r.ID,
		Name:        r.Name,
		Brand:       r.Brand,
		Description: r.Description,
		Price:       r.Price,
		Stock:       r.Stock,
		Category:    r.Category,
		Image:       r.Image,
		IsActive:    r.IsActive,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
