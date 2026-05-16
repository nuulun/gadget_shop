package service

import (
	"context"
	"log"
	"notification-service/internal/model"
)

type Repository interface {
	Create(context.Context, *model.Notification) error
	List(context.Context) ([]model.Notification, error)
}

type Service struct{ repo Repository }

func New(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Send(ctx context.Context, in model.SendNotificationInput) (model.Notification, error) {
	notifType := in.Type
	if notifType == "" {
		notifType = "email"
	}
	n := &model.Notification{
		Recipient: in.Recipient,
		Type:      notifType,
		Message:   in.Message,
		Status:    "sent",
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return model.Notification{}, err
	}
	log.Printf("[notification] sent %s to %s: %s", notifType, in.Recipient, in.Message)
	return *n, nil
}

func (s *Service) List(ctx context.Context) ([]model.Notification, error) {
	return s.repo.List(ctx)
}
