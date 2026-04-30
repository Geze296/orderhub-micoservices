package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Geze296/orderhub/api-service/internal/cache"
	"github.com/Geze296/orderhub/api-service/internal/config"
	"github.com/Geze296/orderhub/api-service/internal/db"
	"github.com/Geze296/orderhub/api-service/internal/event"
	"github.com/Geze296/orderhub/api-service/internal/http/handler"
	"github.com/Geze296/orderhub/api-service/internal/http/routes"
	"github.com/Geze296/orderhub/api-service/internal/logger"
	"github.com/Geze296/orderhub/api-service/internal/repository"
	"github.com/Geze296/orderhub/api-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config     *config.Config
	Logger     *slog.Logger
	DB         *pgxpool.Pool
	Redis      *redis.Client
	HttpServer *http.Server
}

func New(ctx context.Context) (*App, error) {
	cfg := config.LoadConfig()
	log := logger.New(cfg.AppEnv)

	postgres, err := db.NewPostgres(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("Postgres Error: %v", err)
	}

	redisDB, err := cache.NewRedis(ctx, cfg.RedisAddr, cfg.RedisDB)
	if err != nil {
		postgres.Close()
		return nil, fmt.Errorf("Redis Error: %v", err)
	}

	healthHandler := handler.NewHealthHandler()
	authRepository := repository.NewUserRepo(postgres)
	productRepository := repository.NewProductRepository(postgres)
	productCache := cache.NewProductCache(redisDB)
	productService := service.NewProductService(productRepository, *productCache)
	productHandler := handler.NewProductHandler(productService)
	authService := service.NewAuthService(authRepository, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	orderPublisher := event.NewRedisPublisher(redisDB)

	orderRepository := repository.NewOrderRepository(postgres)
	orderService := service.NewOrderService(postgres, orderRepository, productRepository, productCache, orderPublisher, log)
	orderHandler := handler.NewOrderHandler(orderService)

	router := routes.NewRouter(routes.Dependancy{
		Logger:         log,
		Config:         cfg,
		HealthHandler:  healthHandler,
		AuthHandler:    authHandler,
		ProductHandler: productHandler,
		OrderHandler:   orderHandler,
	})

	server := &http.Server{
		Addr:              cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		Config:     cfg,
		Logger:     log,
		DB:         postgres,
		Redis:      redisDB,
		HttpServer: server,
	}, nil
}

func (a *App) Close() {
	if a.Redis != nil {
		_ = a.Redis.Close()
	}
	if a.DB != nil {
		a.DB.Close()
	}
}
