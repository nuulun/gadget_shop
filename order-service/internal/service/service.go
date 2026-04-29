package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/internal/model"
	"order-service/internal/repository"
	"time"
)

var ErrNotFound = errors.New("order not found")

type Repo interface {
	Create(context.Context, model.Order) (model.Order, error)
	GetByID(context.Context, uint64) (model.Order, error)
	ListByUserID(context.Context, uint64) ([]model.Order, error)
	List(context.Context) ([]model.Order, error)
}

type OrderService struct {
	repo              Repo
	productServiceURL string
	httpClient        *http.Client
}

func New(repo Repo, productServiceURL string) *OrderService {
	return &OrderService{
		repo:              repo,
		productServiceURL: productServiceURL,
		httpClient:        &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *OrderService) Create(ctx context.Context, in model.CreateOrderInput) (model.Order, error) {
	if in.UserID == 0 {
		return model.Order{}, fmt.Errorf("user_id is required")
	}
	if len(in.Items) == 0 {
		return model.Order{}, fmt.Errorf("order must have at least one item")
	}

	var items []model.OrderItem
	total := 0.0

	for _, req := range in.Items {
		product, err := s.fetchProduct(ctx, req.ProductID)
		if err != nil {
			return model.Order{}, fmt.Errorf("product %d: %w", req.ProductID, err)
		}
		if product.Stock < req.Quantity {
			return model.Order{}, fmt.Errorf("insufficient stock for product %d (have %d, need %d)", req.ProductID, product.Stock, req.Quantity)
		}
		lineTotal := product.Price * float64(req.Quantity)
		total += lineTotal
		items = append(items, model.OrderItem{ProductID: req.ProductID, Quantity: req.Quantity, Price: product.Price})
	}

	order := model.Order{
		UserID:     in.UserID,
		Status:     "pending",
		TotalPrice: total,
		Items:      items,
	}
	return s.repo.Create(ctx, order)
}

func (s *OrderService) GetByID(ctx context.Context, id uint64) (model.Order, error) {
	o, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Order{}, ErrNotFound
	}
	return o, err
}

func (s *OrderService) ListByUserID(ctx context.Context, userID uint64) ([]model.Order, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *OrderService) List(ctx context.Context) ([]model.Order, error) {
	return s.repo.List(ctx)
}

func (s *OrderService) fetchProduct(ctx context.Context, id uint64) (model.ProductInfo, error) {
	url := fmt.Sprintf("%s/products/%d", s.productServiceURL, id)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return model.ProductInfo{}, fmt.Errorf("product-service unreachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return model.ProductInfo{}, fmt.Errorf("product not found")
	}
	if resp.StatusCode != http.StatusOK {
		return model.ProductInfo{}, fmt.Errorf("product-service error %d", resp.StatusCode)
	}
	var p model.ProductInfo
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return model.ProductInfo{}, err
	}
	return p, nil
}
