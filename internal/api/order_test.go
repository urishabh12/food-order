package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"order-food/internal/api"
	"order-food/internal/config"
)

func TestPlaceOrder_Success_NoCoupon(t *testing.T) {
	cfg := config.Config{Addr: ":0", BasePath: "/api", APIKey: "apitest", CouponsDir: "."}
	h := &api.Handlers{Coupons: loadTestCoupons(t)}
	r := api.NewRouter(cfg, h)

	body := obj(map[string]any{
		"items": []any{
			map[string]any{"productId": "10", "quantity": 2},
			map[string]any{"productId": "31", "quantity": 1},
		},
	})

	w := doReq(t, r, http.MethodPost, "/api/order", body, map[string]string{"api_key": "apitest"})
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	j := mustBodyJSON(t, w)
	if j["id"] == "" {
		t.Fatalf("expected order id, got %v", j["id"])
	}
}

func TestPlaceOrder_Success_WithValidCoupon(t *testing.T) {
	// loadTestCoupons marks HAPPYHRS as valid (present in 2 files)
	cfg := config.Config{Addr: ":0", BasePath: "/api", APIKey: "apitest", CouponsDir: "."}
	h := &api.Handlers{Coupons: loadTestCoupons(t)}
	r := api.NewRouter(cfg, h)

	body := obj(map[string]any{
		"couponCode": "HAPPYHRS",
		"items": []any{
			map[string]any{"productId": "10", "quantity": 1},
		},
	})

	w := doReq(t, r, http.MethodPost, "/api/order", body, map[string]string{"api_key": "apitest"})
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPlaceOrder_AuthMissing(t *testing.T) {
	cfg := config.Config{Addr: ":0", BasePath: "/api", APIKey: "apitest", CouponsDir: "."}
	h := &api.Handlers{Coupons: loadTestCoupons(t)}
	r := api.NewRouter(cfg, h)

	body := obj(map[string]any{
		"items": []any{
			map[string]any{"productId": "10", "quantity": 1},
		},
	})
	w := doReq(t, r, http.MethodPost, "/api/order", body, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPlaceOrder_InvalidCoupon(t *testing.T) {
	// SUPER100 appears in only one file -> invalid per rule (needs >=2)
	cfg := config.Config{Addr: ":0", BasePath: "/api", APIKey: "apitest", CouponsDir: "."}
	h := &api.Handlers{Coupons: loadTestCoupons(t)}
	r := api.NewRouter(cfg, h)

	body := obj(map[string]any{
		"couponCode": "SUPER100",
		"items": []any{
			map[string]any{"productId": "10", "quantity": 1},
		},
	})

	w := doReq(t, r, http.MethodPost, "/api/order", body, map[string]string{"api_key": "apitest"})
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want 422, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPlaceOrder_ValidationErrors(t *testing.T) {
	cfg := config.Config{Addr: ":0", BasePath: "/api", APIKey: "apitest", CouponsDir: "."}
	h := &api.Handlers{Coupons: loadTestCoupons(t)}
	r := api.NewRouter(cfg, h)

	t.Run("empty items", func(t *testing.T) {
		body := obj(map[string]any{
			"items": []any{},
		})
		w := doReq(t, r, http.MethodPost, "/api/order", body, map[string]string{"api_key": "apitest"})
		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("want 422, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("non-positive quantity", func(t *testing.T) {
		body := obj(map[string]any{
			"items": []any{
				map[string]any{"productId": "10", "quantity": 0},
			},
		})
		w := doReq(t, r, http.MethodPost, "/api/order", body, map[string]string{"api_key": "apitest"})
		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("want 422, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("unknown product", func(t *testing.T) {
		body := obj(map[string]any{
			"items": []any{
				map[string]any{"productId": "999", "quantity": 1},
			},
		})
		w := doReq(t, r, http.MethodPost, "/api/order", body, map[string]string{"api_key": "apitest"})
		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("want 422, got %d: %s", w.Code, w.Body.String())
		}
	})
}

func obj(m map[string]any) []byte {
	b, _ := json.Marshal(m)
	return b
}
