package domain

import "time"

type Product struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	PriceCents float64 `json:"price_cents"`
	Stock int32 `json:"stock"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}
