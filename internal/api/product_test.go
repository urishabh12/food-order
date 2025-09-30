package api_test

import (
	"net/http"
	"testing"

	"order-food/internal/api"
	"order-food/internal/config"
)

func TestListProducts(t *testing.T) {
	cfg := config.Config{Addr: ":0", BasePath: "/api", APIKey: "apitest", CouponsDir: "."}
	h := &api.Handlers{Coupons: loadTestCoupons(t)}
	r := api.NewRouter(cfg, h)

	w := doReq(t, r, http.MethodGet, "/api/product", nil, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	// minimally assert it's an array
	if got := w.Body.Bytes(); len(got) == 0 || got[0] != '[' {
		t.Fatalf("expected JSON array, got: %s", string(got))
	}
}

func TestGetProduct(t *testing.T) {
	cfg := config.Config{Addr: ":0", BasePath: "/api", APIKey: "apitest", CouponsDir: "."}
	h := &api.Handlers{Coupons: loadTestCoupons(t)}
	r := api.NewRouter(cfg, h)

	t.Run("ok", func(t *testing.T) {
		w := doReq(t, r, http.MethodGet, "/api/product/10", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
		}
		j := mustBodyJSON(t, w)
		if j["id"] != "10" {
			t.Fatalf("expected id=10, got %v", j["id"])
		}
	})

	t.Run("bad id (non-numeric path param)", func(t *testing.T) {
		w := doReq(t, r, http.MethodGet, "/api/product/abc", nil, nil)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("want 400, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doReq(t, r, http.MethodGet, "/api/product/999999", nil, nil)
		if w.Code != http.StatusNotFound {
			t.Fatalf("want 404, got %d: %s", w.Code, w.Body.String())
		}
	})
}
