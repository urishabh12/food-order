package api

import (
	"net/http"
	"order-food/internal/data"
	"order-food/internal/middleware"
	"order-food/internal/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handlers struct{}

func (h *Handlers) Healthcheck(w http.ResponseWriter, _ *http.Request) {
	health := models.Health{Status: true}
	middleware.WriteJSON(w, http.StatusOK, health)
}

func (h *Handlers) GetProducts(w http.ResponseWriter, _ *http.Request) {
	products := data.AllProducts()
	middleware.WriteJSON(w, http.StatusOK, products)
}

func (h *Handlers) GetProductById(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "productId")
	if productId == "" {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid ID supplied")
		return
	}
	if _, err := strconv.ParseInt(productId, 0, 64); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid ID supplied")
		return
	}
	if p, ok := data.ProductByID(productId); ok {
		middleware.WriteJSON(w, http.StatusOK, p)
		return
	}
	middleware.WriteError(w, http.StatusNotFound, "Product not found")
}
