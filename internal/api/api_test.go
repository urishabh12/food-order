package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"order-food/internal/coupons"
)

// --- test helpers ---

type respJSON map[string]any

func doReq(t *testing.T, r http.Handler, method, path string, body []byte, hdr map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func mustBodyJSON(t *testing.T, w *httptest.ResponseRecorder) respJSON {
	t.Helper()
	var v respJSON
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Fatalf("response not JSON: %v\nbody=%s", err, w.Body.String())
	}
	return v
}

// write three small plain-text coupon files; returns coupons.Store
func loadTestCoupons(t *testing.T) *coupons.Store {
	t.Helper()

	tmp := t.TempDir()

	// file A: HAPPYHRS, FIFTYOFF, SUPER100
	writeFile(t, tmp, "a.txt", `
HAPPYHRS
FIFTYOFF
SUPER100
`)

	// file B: HAPPYHRS only
	writeFile(t, tmp, "b.txt", `
xyz
HAPPYHRS
`)

	// file C: unrelated codes
	writeFile(t, tmp, "c.txt", `
QYYSPC46
NSZ0VMH4
FXX621AO
KVM99X2Q
`)

	store, err := coupons.Load(tmp, "a.txt", "b.txt", "c.txt")
	if err != nil {
		t.Fatalf("load coupons: %v", err)
	}
	return store
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := osWriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

// tiny indirection so go vet/go test are happy in restricted sandboxes
var osWriteFile = func(name string, data []byte, perm uint32) error {
	return writeFileImpl(name, data, perm)
}

func writeFileImpl(name string, data []byte, perm uint32) error {
	return os.WriteFile(name, data, os.FileMode(perm))
}
