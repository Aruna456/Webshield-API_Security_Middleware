// Package main starts the WebShield server.
// It loads the configuration, sets up the router with middleware, and defines API endpoints.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Aruna456/webshield/internal/config"
	"github.com/Aruna456/webshield/internal/handlers"
	"github.com/Aruna456/webshield/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Load configuration from the YAML file.
	// This config determines which middlewares are enabled and their settings.
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the HTTP router using chi.
	r := chi.NewRouter()

	// Apply the composable middleware chain based on the config.
	// Middlewares are applied in order: logging -> rate limit -> api key -> jwt.
	r.Use(middleware.Compose(cfg)...)

	// Define API endpoints.
	// /health: A simple health check, no auth required (but middleware may apply if enabled globally).
	// /public: A public endpoint, perhaps requiring API key.
	// /protected: A protected endpoint, requiring full auth (API key + JWT).
	r.Get("/health", handlers.HealthHandler)
	r.Get("/public", handlers.PublicHandler)
	r.Get("/protected", handlers.ProtectedHandler)

	// Start the HTTP server.
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting WebShield server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
