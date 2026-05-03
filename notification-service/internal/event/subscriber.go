package event

import (
	"context"
	"errors"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	logger   *slog.Logger
	rdb      *redis.Client
	handlers *Handlers
}

func NewSubscriber(
	logger *slog.Logger,
	rdb *redis.Client,
	handlers *Handlers,
) *Subscriber {
	return &Subscriber{
		logger:   logger,
		rdb:      rdb,
		handlers: handlers,
	}
}

func (s *Subscriber) Start(ctx context.Context) error {
	pubsub := s.rdb.Subscribe(ctx, OrderCreatedChannel)
	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		return err
	}

	s.logger.Info("subscriber started", slog.String("channel", OrderCreatedChannel))

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("subscriber shutting down")
			return nil

		case msg, ok := <-ch:
			if !ok {
				return errors.New("pubsub channel closed")
			}

			switch msg.Channel {
			case OrderCreatedChannel:
				if err := s.handlers.HandleOrderCreated(ctx, msg.Payload); err != nil {
					s.logger.Error("failed to handle order.created",
						slog.Any("error", err),
						slog.String("payload", msg.Payload),
					)
				}
			}
		}
	}
}