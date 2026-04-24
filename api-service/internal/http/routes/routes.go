package routes

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Geze296/orderhub/api-service/internal/http/handler"
	appmw "github.com/Geze296/orderhub/api-service/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Dependancy struct {
	Logger        *slog.Logger
	HealthHandler *handler.HealthHandler
	AuthHandler   *handler.AuthHandler
}

func NewRouter(deps Dependancy) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(appmw.RequestLogger(deps.Logger))

	r.Get("/health", deps.HealthHandler.Health)

	r.Post("/register", deps.AuthHandler.Register)
	r.Post("/login", deps.AuthHandler.Login)

	return r
}
