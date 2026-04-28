package domain

import "time"

type Order struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Status string `json:"status"`
	TotalAmountCents float64 `json:"total_amount_cents"`
	Items []OrderItem `json:"items"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItem struct {
	ID int64 `json:"id"`
	OrderID int64 `json:"order_id"`
	ProductID int64 `json:"product_id"`
	Quantity int32 `json:"quantity"`
	UnitPriceCents float64 `json:"unit_price_cents"`
}