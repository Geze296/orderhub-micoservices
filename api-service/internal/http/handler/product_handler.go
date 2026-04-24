package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Geze296/orderhub/api-service/internal/http/helper"
	"github.com/Geze296/orderhub/api-service/internal/service"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

type ProductCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int64  `json:"price_cents"`
	Stock       int32  `json:"stock"`
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req ProductCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Error data decoding")
		return
	}

	err := h.ProductService.Create(r.Context(), service.CreateProductRequest{
		Name: req.Name,
		Description: req.Description,
		PriceCents: req.PriceCents,
		Stock: req.Stock,
	})
	
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.WriteJson(w, nil, http.StatusCreated, "Product Created successfully")
}
