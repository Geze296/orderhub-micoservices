package repository

import (
	"context"
	"fmt"

	"github.com/Geze296/orderhub/api-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepositoryInterface interface {
	Create(ctx context.Context, product *domain.Product) error
	List(ctx context.Context) ([]domain.Product, error)
	GetById(ctx context.Context, id int) (*domain.Product, error)
}

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (r *ProductRepo) Create(ctx context.Context, product *domain.Product) error {

	q := `INSERT INTO products (name, description, price_cents, stock)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at
			`
	err := r.db.QueryRow(ctx, q, product.Name, product.Description, product.PriceCents, product.Stock).Scan(&product.ID, &product.CreatedAT)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) List(ctx context.Context) ([]domain.Product, error) {
	q := `SELECT id, name, description, product_cents, stock, created_at, updated_at FROM products
			ORDER BY id DESC
		`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("List product:%w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0)

	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.PriceCents,
			&product.Stock,
			&product.CreatedAT,
			&product.UpdatedAT,
		)
		if err != nil {
			return nil, fmt.Errorf("Scanning product Error:%w", err)
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepo) GetById(ctx context.Context, id int) (*domain.Product, error) {
	q := `SELECT id, name, description, product_cents, stock, created_at, updated_at FROM products
			WHERE id = $1`

	var product domain.Product
	err := r.db.QueryRow(ctx, q, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.PriceCents,
		&product.Stock,
		&product.CreatedAT,
		&product.UpdatedAT,
	)

	if err != nil {
		return nil, fmt.Errorf("Fetch product by id Error: %w", err)
	}
	return &product, nil
}

func (r *ProductRepo) Update(ctx context.Context, product domain.Product) (*domain.Product, error) {
	q := `UPDATE products SET
			name = $2, 
			description = $3,
			price_cents = $4, 
			stock = $5, 
			updated_at = NOW()
		WHERE id = $1 
		`

	err := r.db.QueryRow(ctx, q, product.ID, product.Name, product.Description, product.PriceCents, product.Stock).Scan(
		&product.Name,
		&product.Description,
		&product.PriceCents,
		&product.Stock,
	)

	if err != nil {
		return nil, fmt.Errorf("Error :%w", err)
	}

	return &product, nil
}
