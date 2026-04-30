package event

import "time"

const OrderCreatedChannel = "order.created"

type OrderCreatedEvent struct {
	EventType        string    `json:"event_type"`
	OrderID          int64     `json:"order_id"`
	UserID           int64     `json:"user_id"`
	TotalAmountCents float64   `json:"total_amout_cents"`
	CreatedAt        time.Time `json:"created_at"`
}
