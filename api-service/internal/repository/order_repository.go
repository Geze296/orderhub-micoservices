package repository

import (
	"context"
	"fmt"

	"github.com/Geze296/orderhub/api-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepositoryInterface interface {
	Create(ctx context.Context, q DBTX, order *domain.Order) error
	CreateItem(ctx context.Context, q DBTX, item *domain.OrderItem) error
	ListByUserID(ctx context.Context, userID int64) ([]domain.Order, error)
	GetByOrderIDAndUserID(ctx context.Context, orderID, userID int64) (*domain.Order, error)
}

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, q DBTX, order *domain.Order) error {
	const qy = `INSERT INTO orders (user_id, status, total_amount_cents)
				VALUES ($1, $2, $3)
				RETURNING id, created_at`
				
	err := q.QueryRow(ctx, qy, 
		order.UserID,
		order.Status,
		order.TotalAmountCents,
	).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return fmt.Errorf("Error While creating order")
	}
	return nil
}

func (r *OrderRepository) CreateItem(ctx context.Context, q DBTX, item *domain.OrderItem) error {
	qy := `INSERT INTO order_items 
			(order_id, product_id, quantity, unit_price_cents)
			VALUES ($1, $2, $3, $4)
			RETURNING id`
	err := q.QueryRow(ctx, qy, 
			item.OrderID,
			item.ProductID,
			item.Quantity,
			item.UnitPriceCents,
		).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("Error while creating item")
	}

	return nil
}