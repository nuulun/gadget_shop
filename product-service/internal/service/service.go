package service

import (
	"context"
	"errors"
	"fmt"

	"product-service/internal/model"
	"product-service/internal/repository"
)

// ProductRepository is the interface the service depends on.
type ProductRepository interface {
	List(ctx context.Context, f model.ListFilter) ([]model.Product, error)
	Count(ctx context.Context, f model.ListFilter) (int64, error)
	GetByID(ctx context.Context, id uint64) (model.Product, error)
	Create(ctx context.Context, in model.CreateProductInput) (model.Product, error)
	Update(ctx context.Context, id uint64, in model.UpdateProductInput) (model.Product, error)
	Delete(ctx context.Context, id uint64) error
}

// ErrNotFound is the service-level not-found error.
var ErrNotFound = errors.New("product not found")

// ListResult wraps a page of products with pagination metadata.
type ListResult struct {
	Items  []model.Product `json:"items"`
	Total  int64           `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// ProductService contains all product business logic.
type ProductService struct {
	repo ProductRepository
}

// New creates a ProductService.
func New(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// List returns a paginated, optionally filtered list of products.
func (s *ProductService) List(ctx context.Context, f model.ListFilter) (ListResult, error) {
	products, err := s.repo.List(ctx, f)
	if err != nil {
		return ListResult{}, fmt.Errorf("service.List: %w", err)
	}
	total, err := s.repo.Count(ctx, f)
	if err != nil {
		return ListResult{}, fmt.Errorf("service.List count: %w", err)
	}
	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	return ListResult{Items: products, Total: total, Limit: limit, Offset: f.Offset}, nil
}

// GetByID retrieves a single product.
func (s *ProductService) GetByID(ctx context.Context, id uint64) (model.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Product{}, ErrNotFound
	}
	if err != nil {
		return model.Product{}, fmt.Errorf("service.GetByID: %w", err)
	}
	return p, nil
}

// Create validates and creates a new product.
func (s *ProductService) Create(ctx context.Context, in model.CreateProductInput) (model.Product, error) {
	if in.Name == "" {
		return model.Product{}, fmt.Errorf("name is required")
	}
	if in.Price < 0 {
		return model.Product{}, fmt.Errorf("price must be non-negative")
	}
	p, err := s.repo.Create(ctx, in)
	if err != nil {
		return model.Product{}, fmt.Errorf("service.Create: %w", err)
	}
	return p, nil
}

// Update applies partial updates to a product.
func (s *ProductService) Update(ctx context.Context, id uint64, in model.UpdateProductInput) (model.Product, error) {
	if in.Price != nil && *in.Price < 0 {
		return model.Product{}, fmt.Errorf("price must be non-negative")
	}
	p, err := s.repo.Update(ctx, id, in)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Product{}, ErrNotFound
	}
	if err != nil {
		return model.Product{}, fmt.Errorf("service.Update: %w", err)
	}
	return p, nil
}

// Delete soft-deletes a product.
func (s *ProductService) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
