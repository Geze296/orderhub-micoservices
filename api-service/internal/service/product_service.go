package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Geze296/orderhub/api-service/internal/domain"
	"github.com/Geze296/orderhub/api-service/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepositoryInterface
}

func NewProductService(repo repository.ProductRepositoryInterface) *ProductService {
	return &ProductService{repo: repo}
}


var (
	ErrInvalidProductInput = errors.New("Invalid product input")
)


type CreateProductRequest struct {
	Name        string
	Description string
	PriceCents  int64
	Stock       int32
}


func (s *ProductService) Create(ctx context.Context, input CreateProductRequest) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" || input.Description == "" || input.PriceCents <= 0 || input.Stock <= 0 {
		return ErrInvalidProductInput
	}

	product := domain.Product{
		Name:        input.Name,
		Description: input.Description,
		PriceCents:  input.PriceCents,
		Stock:       input.Stock,
	}

	err := s.repo.Create(ctx, &product)
	if err != nil {
		return err
	}
	return nil
}
