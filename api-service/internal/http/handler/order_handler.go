package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Geze296/orderhub/api-service/internal/http/helper"
	"github.com/Geze296/orderhub/api-service/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type CreateOrderReq struct {
	Items []service.CreateOrderItem `json:"items"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJson(w, nil, http.StatusOK, "Correctly order created")
}
