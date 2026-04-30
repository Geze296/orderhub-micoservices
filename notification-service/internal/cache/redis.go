package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, addr string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB: db,
		DialTimeout: 5*time.Second,
		ReadTimeout: 5*time.Second,
		WriteTimeout: 5*time.Second,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Redis ping error")
	}
	fmt.Println("Redis connected!")
	return client, nil
}