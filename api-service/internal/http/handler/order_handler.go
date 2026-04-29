package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Geze296/orderhub/api-service/internal/http/helper"
	"github.com/Geze296/orderhub/api-service/internal/http/middleware"
	"github.com/Geze296/orderhub/api-service/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type CreateOrderReq struct {
	Items []CreateOrderItemReq `json:"items"`
}

type CreateOrderItemReq struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.UserIdFromContext(r.Context())
	if !ok {
		helper.WriteError(w, http.StatusBadRequest, "Error in User id context")
		return
	}

	var req CreateOrderReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	items := make([]service.CreateOrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		items = append(items, service.CreateOrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order, err := h.orderService.Create(r.Context(), service.CreateOrderInput{
		UserID:     userId,
		OrderItems: items,
	})
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJson(w, order, http.StatusCreated, "Correctly order created")
}

func (h *OrderHandler) ListUserOrders(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.UserIdFromContext(r.Context())
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "Wrong userid")
		return
	}
	orders, err := h.orderService.ListByUserID(r.Context(), userId)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJson(w, orders, http.StatusOK, "User Order fetched successfully")
}

func (h *OrderHandler) ListByUserANDOrderID(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.UserIdFromContext(r.Context())
	orderID, err := parseIntParams(r, "id")
	
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "Wrong userid")
		return
	}
	orders, e := h.orderService.ListByUserANDOrderID(r.Context(), int64(orderID), userId)
	if e != nil {
		helper.WriteError(w, http.StatusInternalServerError, e.Error())
		return
	}

	helper.WriteJson(w, orders, http.StatusOK, "User Order fetched successfully")
}
