package service

import (
	"context"
	"errors"
	"fmt"
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

type UpdateProductRequest struct {
	ID          int64
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

func (s *ProductService) List(ctx context.Context) ([]domain.Product, error) {
	
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ProductService) GetById(ctx context.Context, id int) (*domain.Product, error) {
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	product, err := s.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Update(ctx context.Context, input UpdateProductRequest) (*domain.Product, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.PriceCents < 0 || input.Stock < 0 {
		return nil, fmt.Errorf("Price and stock should greater than 0")
	}
	product := domain.Product{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		PriceCents:  input.PriceCents,
		Stock:       input.Stock,
	}

	updateProduct, err := s.repo.Update(ctx, &product)
	if err != nil {
		return nil, err
	}
	return updateProduct, err
}

func (s *ProductService) Delete(ctx context.Context, id int) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}