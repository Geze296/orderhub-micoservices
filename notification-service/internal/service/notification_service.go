package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Geze296/orderhub/notification-service/internal/domain"
	"github.com/Geze296/orderhub/notification-service/internal/repository"
)

type NotificationService struct {
	logger        *slog.Logger
	notifications repository.NotificationRepository
}

func NewNotificationService(
	logger *slog.Logger,
	notifications repository.NotificationRepository,
) *NotificationService {
	return &NotificationService{
		logger:        logger,
		notifications: notifications,
	}
}

func (s *NotificationService) HandleOrderCreated(
	ctx context.Context,
	userID int64,
	orderID int64,
	totalAmountCents int64,
) error {
	message := fmt.Sprintf("Your order #%d has been created successfully. Total: %d cents", orderID, totalAmountCents)

	notification := &domain.Notification{
		UserID:  userID,
		OrderID: orderID,
		Type:    "order.created",
		Message: message,
		Status:  "sent",
	}

	if err := s.notifications.Create(ctx, notification); err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	s.logger.Info("notification sent",
		slog.Int64("user_id", userID),
		slog.Int64("order_id", orderID),
		slog.Int64("notification_id", notification.ID),
	)

	return nil
}