package api

import (
	"encoding/json"
	"net/http"
	"order-food/internal/models"
)

type Handlers struct{}

func (h *Handlers) Healthcheck(w http.ResponseWriter, _ *http.Request) {
	health := models.Health{Status: true}
	writeJSON(w, 200, health)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(marshal(v))
}

func marshal(v any) []byte {
	b, _ := jsonMarshal(v)
	return b
}

var jsonMarshal = func(v any) ([]byte, error) { return json.Marshal(v) }
