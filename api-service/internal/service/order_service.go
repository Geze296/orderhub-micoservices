package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Geze296/orderhub/api-service/internal/cache"
	"github.com/Geze296/orderhub/api-service/internal/domain"
	"github.com/Geze296/orderhub/api-service/internal/event"
	"github.com/Geze296/orderhub/api-service/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidInputOrder = errors.New("Invalid Input")
	ErrInsufficientStock = errors.New("Insufficient Stock amount")
)

type OrderService struct {
	db           *pgxpool.Pool
	orderRepo    *repository.OrderRepository
	productRepo  *repository.ProductRepo
	productCache *cache.ProductCache
	publisher    *event.Publisher
	logger       *slog.Logger
}

func NewOrderService(db *pgxpool.Pool, orderRepo *repository.OrderRepository, productRepo *repository.ProductRepo, productCache *cache.ProductCache, publisher *event.Publisher, logger *slog.Logger) *OrderService {
	return &OrderService{
		db:           db,
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		productCache: productCache,
		publisher:    publisher,
		logger:       logger,
	}
}

type CreateOrderInput struct {
	UserID     int64             `json:"_"`
	OrderItems []CreateOrderItem `json:"items"`
}

type CreateOrderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func (s *OrderService) Create(ctx context.Context, input CreateOrderInput) (*domain.Order, error) {
	if input.UserID <= 0 || len(input.OrderItems) <= 0 {
		return nil, ErrInvalidInputOrder
	}

	seen := make(map[int64]bool, len(input.OrderItems))

	for _, item := range input.OrderItems {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return nil, ErrInvalidInputOrder
		}
		if seen[item.ProductID] {
			return nil, ErrInvalidInputOrder
		}
		seen[item.ProductID] = true
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error Beginning ")
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	order := &domain.Order{
		UserID: input.UserID,
		Status: "Pending",
		Items:  make([]domain.OrderItem, 0, len(input.OrderItems)),
	}

	affectedProductsIDs := make([]int64, 0, len(input.OrderItems))

	for _, reqItem := range input.OrderItems {
		product, err := s.productRepo.GetByIdForUpdate(ctx, tx, reqItem.ProductID)
		if err != nil {
			return nil, err
		}
		if reqItem.Quantity > product.Stock {
			return nil, ErrInsufficientStock
		}

		newStock := product.Stock - reqItem.Quantity

		e := s.productRepo.UpdateStock(ctx, tx, reqItem.ProductID, newStock)
		if e != nil {
			return nil, e
		}

		orderItem := domain.OrderItem{
			ProductID:      product.ID,
			Quantity:       reqItem.Quantity,
			UnitPriceCents: product.PriceCents,
		}
		order.Items = append(order.Items, orderItem)
		order.TotalAmountCents += product.PriceCents + float64(orderItem.Quantity)
		affectedProductsIDs = append(affectedProductsIDs, product.ID)
	}

	if err := s.orderRepo.Create(ctx, tx, order); err != nil {
		return nil, err
	}
	fmt.Printf("Repo Order ID: %v", order.ID)

	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		if err := s.orderRepo.CreateItem(ctx, tx, &order.Items[i]); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx error : %w", err)
	}

	_ = s.productCache.DeleteProductList(ctx)
	for _, productID := range affectedProductsIDs {
		s.productCache.DeleteProduct(ctx, int(productID))
	}

	orderEvent := event.OrderCreatedEvent{
		EventType:        event.OrderCreatedChannel,
		OrderID:          order.ID,
		UserID:           order.UserID,
		TotalAmountCents: order.TotalAmountCents,
		CreatedAt:        order.CreatedAt,
	}

	if err := s.publisher.PublishOrderCreated(ctx, orderEvent); err != nil {
		s.logger.Warn("failed to publish order.created event",
			slog.Int64("order_id", order.ID),
			slog.Any("error", err),
		)
	}

	return order, nil
}

func (s *OrderService) ListByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	orders, err := s.orderRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) ListByUserANDOrderID(ctx context.Context, orderID, userID int64) (*domain.Order, error) {
	orders, err := s.orderRepo.GetByOrderIDAndUserID(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	return orders, err
}
