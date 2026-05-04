package dto

type CreateOrderReq struct {
	Items []CreateOrderItemReq `json:"items"`
}

type CreateOrderItemReq struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}
