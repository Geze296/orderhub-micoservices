package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Geze296/orderhub/notification-service/internal/cache"
	"github.com/Geze296/orderhub/notification-service/internal/config"
	"github.com/Geze296/orderhub/notification-service/internal/db"
	"github.com/Geze296/orderhub/notification-service/internal/event"
	"github.com/Geze296/orderhub/notification-service/internal/logger"
	"github.com/Geze296/orderhub/notification-service/internal/repository"
	"github.com/Geze296/orderhub/notification-service/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.AppEnv)

	pg, err := db.NewPostgres(ctx, cfg.PostgresURL)
	if err != nil {
		log.Error("failed to connect postgres", "error", err)
		os.Exit(1)
	}
	defer pg.Close()

	rdb, err := cache.NewRedis(ctx, cfg.RedisAddr, cfg.RedisDB)
	if err != nil {
		log.Error("failed to connect redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()

	notificationRepo := repository.NewPostgresNotificationRepository(pg)
	notificationService := service.NewNotificationService(log, notificationRepo)

	eventHandlers := event.NewHandlers(log, notificationService)
	subscriber := event.NewSubscriber(log, rdb, eventHandlers)

	log.Info("starting notification-service")

	if err := subscriber.Start(ctx); err != nil {
		log.Error("subscriber stopped with error", "error", err)
		os.Exit(1)
	}

	log.Info("notification-service stopped")

}
