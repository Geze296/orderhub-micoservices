package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Geze296/orderhub/notification-service/internal/service"
)

const OrderCreatedChannel = "order.created"

type OrderCreatedEvent struct {
	EventType        string    `json:"event_type"`
	OrderID          int64     `json:"order_id"`
	UserID           int64     `json:"user_id"`
	TotalAmountCents int64     `json:"total_amount_cents"`
	CreatedAt        time.Time `json:"created_at"`
}

type Handlers struct {
	logger              *slog.Logger
	notificationService *service.NotificationService
}

func NewHandlers(
	logger *slog.Logger,
	notificationService *service.NotificationService,
) *Handlers {
	return &Handlers{
		logger:              logger,
		notificationService: notificationService,
	}
}

func (h *Handlers) HandleOrderCreated(ctx context.Context, payload string) error {
	var event OrderCreatedEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		return fmt.Errorf("unmarshal order.created event: %w", err)
	}

	h.logger.Info("received order.created event",
		slog.Int64("order_id", event.OrderID),
		slog.Int64("user_id", event.UserID),
		slog.Int64("total_amount_cents", event.TotalAmountCents),
	)

	if err := h.notificationService.HandleOrderCreated(
		ctx,
		event.UserID,
		event.OrderID,
		event.TotalAmountCents,
	); err != nil {
		return fmt.Errorf("handle order.created: %w", err)
	}

	return nil
}