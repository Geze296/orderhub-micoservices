package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Geze296/orderhub/api-service/internal/app"
)

func main() {
	ctx := context.Background()
	application, err := app.New(ctx)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	defer application.Close()

	go func() {
		application.Logger.Info("starting api-service", "port", application.Config.HTTPPort)
		if err := application.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			application.Logger.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	application.Logger.Info("shutting down api-service")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.HttpServer.Shutdown(shutdownCtx); err != nil {
		application.Logger.Error("graceful shutdown failed", "error", err)
	}
}
