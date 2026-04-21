package routes

import (
	"encoding/json"
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
}

func NewRouter(deps Dependancy) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(appmw.RequestLogger(deps.Logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "Ok",
			"message": "Success",
		})
	})
	r.Get("/health", deps.HealthHandler.Health)

	return r
}
