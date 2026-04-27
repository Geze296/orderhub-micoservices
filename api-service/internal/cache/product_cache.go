package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Geze296/orderhub/api-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

type ProductCache struct {
	rdb       *redis.Client
	detailTTL time.Duration
	listTTL   time.Duration
}

func NewProductCache(rdb *redis.Client) *ProductCache {
	return &ProductCache{
		rdb:       rdb,
		detailTTL: 5 * time.Minute,
		listTTL:   1 * time.Minute,
	}
}

func (c *ProductCache) GetProduct(ctx context.Context, id int64) (*domain.Product, bool, error) {
	key := fmt.Sprintf("product:%d", id)

	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var product domain.Product

	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, false, err
	}

	return &product, true, nil
}

func (c *ProductCache) SetProduct(ctx context.Context, product *domain.Product) error {
	key := fmt.Sprintf("product:%d", product.ID)

	b, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, b, c.detailTTL).Err()
}

func (c *ProductCache) DeleteProduct(ctx context.Context, id int) error {
	key := fmt.Sprintf("product:%d", id)

	return c.rdb.Del(ctx, key).Err()
}

func (c *ProductCache) GetProductsList(ctx context.Context) ([]domain.Product, bool, error) {
	key := "product:list"

	val, err := c.rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var productList []domain.Product

	if err := json.Unmarshal([]byte(val), &productList); err != nil {
		return nil, false, err
	}

	return productList, true, nil
}

func (c *ProductCache) SetProductList(ctx context.Context, products []domain.Product) error {
	key := "product:list"

	b, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, b, c.listTTL).Err()
}

func (c *ProductCache) DeleteProductList(ctx context.Context) error {
	key := "product:list"

	return c.rdb.Del(ctx, key).Err()
}
