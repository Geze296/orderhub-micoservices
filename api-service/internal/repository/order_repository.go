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
	fmt.Printf("Order ID: %v", item.OrderID)
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
		return fmt.Errorf("Error while creating item: %w", err)
	}

	return nil
}

func (r *OrderRepository) ListByUserID(ctx context.Context, userId int64) ([]domain.Order, error) {
	qy := `SELECT id, user_id, status, total_amount_cents, created_at
			FROM orders
			WHERE user_id = $1
			ORDER BY id DESC
		`
	rows, err := r.db.Query(ctx, qy, userId)
	if err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)

	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmountCents,
			&order.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("Error: %w", err)
		}

		items, err := r.listOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		order.Items = items
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) listOrderItemsByOrderID(ctx context.Context, order_id int64) ([]domain.OrderItem, error) {
	qy := `SELECT id, order_id, product_id, quantity, unit_price_cents
			FROM order_items
			WHERE order_id = $1
			ORDER BY id ASC
		`
	rows, err := r.db.Query(ctx, qy, order_id)
	if err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}
	defer rows.Close()

	order_items := make([]domain.OrderItem, 0)
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPriceCents,
		); err != nil {
			return nil, fmt.Errorf("Error: %w", err)
		}
		order_items = append(order_items, item)
	}

	return order_items, nil
}

func (r *OrderRepository) GetByOrderIDAndUserID(ctx context.Context, orderID, userID int64) (*domain.Order, error) {
	qy := `SELECT id, user_id, status, total_amount_cents, created_at
			FROM orders
			WHERE id = $1 AND user_id = $2
			ORDER BY id DESC
		`
	var order domain.Order
	if err := r.db.QueryRow(ctx, qy, orderID, userID).Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmountCents,
			&order.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}
	items, err := r.listOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return &order, nil
}
