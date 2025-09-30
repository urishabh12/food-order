package middleware

import (
	"encoding/json"
	"net/http"
)

func RequireAPIKey(expected string) func(http.Handler) http.Handler {
	const headerName = "api_key"
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Header.Get(headerName)
			if expected == "" || got == expected {
				next.ServeHTTP(w, r)
				return
			}
			if got == "" {
				WriteError(w, http.StatusUnauthorized, "missing api_key header")
				return
			}
			WriteError(w, http.StatusUnauthorized, "invalid api_key")
		})
	}
}

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
