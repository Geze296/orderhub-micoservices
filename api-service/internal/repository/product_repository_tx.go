package repository

import (
	"context"
	"fmt"

	"github.com/Geze296/orderhub/api-service/internal/domain"
)

func (r *ProductRepo) GetByIdForUpdate(ctx context.Context, q DBTX, id int64) (*domain.Product, error) {
	const qy = `
		SELECT id, name, description, price_cents, stock, created_at, updated_at
		FROM products
		WHERE id = $1
		FOR UPDATE
	`

	var p domain.Product

	err := q.QueryRow(ctx, qy, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.PriceCents,
		&p.Stock,
		&p.CreatedAT,
		&p.UpdatedAT,
	)
	if err != nil {
		return nil, fmt.Errorf("Get product for update error: %w", err)
	}

	return &p, nil
}

func (r *ProductRepo) UpdateStock(ctx context.Context, q DBTX, id int64, stock int32) error {
	qy := `UPDATE products
			SET stock = $2, updated_at = NOW()
			WHERE id = $1
		`
	tag, err := q.Exec(ctx, qy, id, stock)
	if err != nil {
		return fmt.Errorf("Error while updating stock")
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("There is no table is affected")
	}
	
	return nil
}