package middleware

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(mustJSON(v))
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

func mustJSON(v any) []byte {
	b, _ := jsonMarshal(v)
	return b
}

// tiny wrapper to allow tree-shaking if replaced later
var jsonMarshal = func(v any) ([]byte, error) { return json.Marshal(v) }
