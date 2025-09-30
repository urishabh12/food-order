package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-food/internal/coupons"
	"order-food/internal/data"
	"order-food/internal/middleware"
	"order-food/internal/models"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handlers struct {
	Coupons *coupons.Store
}

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

func (h *Handlers) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req models.OrderReq
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if req.CouponCode != nil && strings.TrimSpace(*req.CouponCode) != "" {
		if ok, reason := h.Coupons.Valid(*req.CouponCode); !ok {
			middleware.WriteError(w, http.StatusUnprocessableEntity, "Validation exception: "+reason)
			return
		}
	}

	if len(req.Items) == 0 {
		middleware.WriteError(w, http.StatusUnprocessableEntity, "Validation exception: items required")
		return
	}
	unique := map[string]models.Product{}
	for i, it := range req.Items {
		if strings.TrimSpace(it.ProductID) == "" {
			middleware.WriteError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Validation exception: items[%d].productId required", i))
			return
		}
		if it.Quantity <= 0 {
			middleware.WriteError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Validation exception: items[%d].quantity must be > 0", i))
			return
		}
		p, ok := data.ProductByID(it.ProductID)
		if !ok {
			middleware.WriteError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Validation exception: product %q not found", it.ProductID))
			return
		}
		unique[p.ID] = p
	}

	order := models.Order{
		ID:    uuid.NewString(),
		Items: req.Items,
	}
	for _, p := range unique {
		order.Products = append(order.Products, p)
	}
	middleware.WriteJSON(w, http.StatusOK, order)
}
