package api_test

import (
	"net/http"
	"testing"

	"order-food/internal/api"
	"order-food/internal/config"
)

func TestHealthz(t *testing.T) {
	cfg := config.Config{Addr: ":0", BasePath: "/api", APIKey: "apitest", CouponsDir: "."}
	h := &api.Handlers{Coupons: loadTestCoupons(t)}
	r := api.NewRouter(cfg, h)

	w := doReq(t, r, http.MethodGet, "/api/healthz", nil, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
}
