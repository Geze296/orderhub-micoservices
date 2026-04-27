package routes

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Geze296/orderhub/api-service/internal/config"
	"github.com/Geze296/orderhub/api-service/internal/http/handler"
	appmw "github.com/Geze296/orderhub/api-service/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Dependancy struct {
	Logger         *slog.Logger
	Config         *config.Config
	HealthHandler  *handler.HealthHandler
	AuthHandler    *handler.AuthHandler
	ProductHandler *handler.ProductHandler
}

func NewRouter(deps Dependancy) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(appmw.RequestLogger(deps.Logger))

	r.Get("/health", deps.HealthHandler.Health)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", deps.AuthHandler.Register)
			r.Post("/login", deps.AuthHandler.Login)
		})
		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware(deps.Config.JWTSecret))
			r.Get("/me", deps.AuthHandler.GetMe)
		})
		r.Route("/product", func(r chi.Router) {
			r.Use(appmw.AuthMiddleware(deps.Config.JWTSecret))
			r.Post("/create", deps.ProductHandler.CreateProduct)
			r.Get("/", deps.ProductHandler.GetAllProducts)
			r.Get("/{id}", deps.ProductHandler.GetById)
			r.Put("/{id}", deps.ProductHandler.UpdateProduct)
			r.Delete("/{id}", deps.ProductHandler.DeleteProduct)
		})
	})

	return r
}
