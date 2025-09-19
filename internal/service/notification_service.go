package service

import (
	"context"
	"log/slog"

	"notification-service/repo"
    "notification-service/entity"
)

type NotificationService struct {
	repo *repo.NotificationRepo
}

func NewNotificationService(r *repo.NotificationRepo) *NotificationService {
	return &NotificationService{repo: r}
}

func (s *NotificationService) SaveNotification(ctx context.Context, n *model.Notification) (*model.Notification, error) {
	if err := s.repo.Save(ctx, n); err != nil {
		slog.Error("save notification", "error", err)
		return nil, err
	}
	slog.Info("saved notification", "id", n.ID)
	return n, nil
}

