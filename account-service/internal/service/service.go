package service

import (
	"account-service/internal/model"
	"account-service/internal/repository"
	"context"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("user not found")

type Repo interface {
	Create(context.Context, model.CreateUserInput) (model.User, error)
	GetByID(context.Context, uint64) (model.User, error)
	List(context.Context, int, int) ([]model.User, error)
	Update(context.Context, uint64, model.UpdateUserInput) (model.User, error)
	Delete(context.Context, uint64) error
}

type AccountService struct{ repo Repo }

func New(repo Repo) *AccountService { return &AccountService{repo: repo} }

func (s *AccountService) Create(ctx context.Context, in model.CreateUserInput) (model.User, error) {
	if in.Login == "" {
		return model.User{}, fmt.Errorf("login is required")
	}
	return s.repo.Create(ctx, in)
}

func (s *AccountService) GetByID(ctx context.Context, id uint64) (model.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return model.User{}, ErrNotFound
	}
	return u, err
}

func (s *AccountService) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	if limit <= 0 { limit = 20 }
	return s.repo.List(ctx, limit, offset)
}

func (s *AccountService) Update(ctx context.Context, id uint64, in model.UpdateUserInput) (model.User, error) {
	u, err := s.repo.Update(ctx, id, in)
	if errors.Is(err, repository.ErrNotFound) {
		return model.User{}, ErrNotFound
	}
	return u, err
}

func (s *AccountService) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
