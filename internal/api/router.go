package api

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"order-food/internal/config"
	"order-food/internal/middleware"
)

func NewRouter(cfg config.Config, h *Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Logger)

	r.Route(cfg.BasePath, func(api chi.Router) {
		api.Get("/healthz", h.Healthcheck)
		api.Get("/product", h.GetProducts)
		api.Get("/product/{productId}", h.GetProductById)
		api.With(middleware.RequireAPIKey(cfg.APIKey)).Post("/order", h.PlaceOrder)
	})

	return r
}
