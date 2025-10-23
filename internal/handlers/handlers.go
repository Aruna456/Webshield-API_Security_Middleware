// Package handlers defines the API endpoint handlers.
// These are simple functions that respond to requests after middleware checks.
package handlers

import (
	"net/http"
)

// HealthHandler responds with a simple OK for health checks.
// Useful for monitoring if the server is up.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is healthy"))
}

// PublicHandler is a public endpoint example.
// It can be accessed with minimal auth, depending on config.
func PublicHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is a public resource"))
}

// ProtectedHandler is a protected endpoint example.
// It requires full authentication as per middleware.
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is a protected resource"))
}
