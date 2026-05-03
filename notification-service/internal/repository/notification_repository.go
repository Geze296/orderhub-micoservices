package repository

import (
	"context"
	"fmt"

	"github.com/Geze296/orderhub/notification-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
}

type PostgresNotificationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresNotificationRepository(db *pgxpool.Pool) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{db: db}
}

func (r *PostgresNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	const q = `
		INSERT INTO notifications (user_id, order_id, type, message, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, q,
		notification.UserID,
		notification.OrderID,
		notification.Type,
		notification.Message,
		notification.Status,
	).Scan(&notification.ID, &notification.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}

	return nil
}