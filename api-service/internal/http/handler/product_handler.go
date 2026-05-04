package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Geze296/orderhub/api-service/internal/http/dto"
	"github.com/Geze296/orderhub/api-service/internal/http/helper"
	"github.com/Geze296/orderhub/api-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}



func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Error data decoding")
		return
	}

	err := h.ProductService.Create(r.Context(), service.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		PriceCents:  req.PriceCents,
		Stock:       req.Stock,
	})

	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.WriteJson(w, nil, http.StatusCreated, "Product Created successfully")
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductService.List(r.Context())

	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.WriteJson(w, products, http.StatusOK, "Products fetched successfully")
}

func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParams(r, "id")
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "Error in id parsing")
		return
	}
	product, err := h.ProductService.GetById(r.Context(), id)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJson(w, product, http.StatusOK, "Product fetch successfully")
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParams(r, "id")
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "Error parsing id params")
		return
	}

	if _, err := h.ProductService.GetById(r.Context(), id); err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.ProductCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Invalid input while decoding json")
		return
	}

	updatedProduct, err := h.ProductService.Update(r.Context(), service.UpdateProductRequest{
		ID:          int64(id),
		Name:        req.Name,
		Description: req.Description,
		PriceCents:  req.PriceCents,
		Stock:       req.Stock,
	})

	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.WriteJson(w, updatedProduct, http.StatusOK, "Product updated Successfully")
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseIntParams(r, "id")
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.ProductService.Delete(r.Context(), id); err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJson(w, nil, http.StatusOK, "Product deleted successfully")
}

func parseIntParams(r *http.Request, name string) (int, error) {
	raw := chi.URLParam(r, name)
	return strconv.Atoi(raw)
}
