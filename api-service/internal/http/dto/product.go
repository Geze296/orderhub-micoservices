package dto

type ProductCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	PriceCents  float64 `json:"price_cents"`
	Stock       int32   `json:"stock"`
}
