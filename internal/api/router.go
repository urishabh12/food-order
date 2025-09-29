package api

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"order-food/internal/config"
)

func NewRouter(cfg config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Logger)

	// API routes
	h := &Handlers{}
	r.Route(cfg.BasePath, func(api chi.Router) {
		api.Get("/healthz", h.Healthcheck)
		api.Get("/product", h.GetProducts)
		api.Get("/product/{productId}", h.GetProductById)
	})

	return r
}
