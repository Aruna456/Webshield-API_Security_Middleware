// Package middleware contains security middlewares.
// This file implements API key authentication.
package middleware

import (
	"net/http"

	"github.com/Aruna456/webshield/internal/config"
)

// APIKey returns a middleware function that checks for a valid API key in the header.
// If the key doesn't match, it returns 401 Unauthorized.
func APIKey(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-Api-Key")
			if apiKey != cfg.APIKey {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}
			// If valid, pass to the next handler.
			next.ServeHTTP(w, r)
		})
	}
}
