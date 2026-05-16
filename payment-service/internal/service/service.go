package service

import (
	"context"
	"errors"
	"math/rand"
	"payment-service/internal/model"
	"time"
)

var ErrNotFound = errors.New("payment not found")

type Repository interface {
	Create(context.Context, *model.Payment) error
	GetByID(context.Context, uint64) (*model.Payment, error)
	GetByOrderID(context.Context, uint64) (*model.Payment, error)
	List(context.Context) ([]model.Payment, error)
}

type Service struct{ repo Repository }

func New(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Process(ctx context.Context, in model.CreatePaymentInput) (model.Payment, error) {
	method := in.Method
	if method == "" {
		method = "card"
	}
	// simulate payment: 95% success rate
	status := "success"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Intn(100) < 5 {
		status = "failed"
	}
	p := &model.Payment{
		OrderID: in.OrderID,
		Amount:  in.Amount,
		Method:  method,
		Status:  status,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return model.Payment{}, err
	}
	return *p, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (model.Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Payment{}, ErrNotFound
	}
	return *p, nil
}

func (s *Service) GetByOrderID(ctx context.Context, orderID uint64) (model.Payment, error) {
	p, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		return model.Payment{}, ErrNotFound
	}
	return *p, nil
}

func (s *Service) List(ctx context.Context) ([]model.Payment, error) {
	return s.repo.List(ctx)
}
