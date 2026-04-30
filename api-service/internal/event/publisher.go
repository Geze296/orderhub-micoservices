package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	rdb *redis.Client
}

func NewRedisPublisher(rdb *redis.Client) *Publisher{
	return &Publisher{rdb: rdb}
}

func (p *Publisher) PublishOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	e := p.rdb.Publish(ctx, OrderCreatedChannel, payload)
	if e != nil {
		return e.Err()
	}
	return nil
}