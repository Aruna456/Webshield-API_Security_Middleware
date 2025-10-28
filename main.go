package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Aruna456/webshield/handlers"
	"github.com/Aruna456/webshield/middleware"
)

func main() {
	// Creating a  Router
	mux := http.NewServeMux()
	handler := middleware.JSONLoggingMiddleware( //  LOGS EVERYTHING
		middleware.SecureHeadersMiddleware( //  HEADERS
			middleware.RateLimitMiddleware(
				middleware.JWTMiddleware(
					middleware.SanitizeMiddleware(
						http.HandlerFunc(handlers.UsersHandler),
						middleware.WithQuery(),
						middleware.WithBody(),
						middleware.WithAllowedFields(map[string][]string{
							"name":  {"string"},
							"age":   {"numeric"},
							"email": {"string"},
						}),
					),
				),
				10, time.Minute,
			),
		),
	)
	//Register endpoints
	mux.Handle("/api/users", handler)
	// start server
	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed:%v", err)
	}
}
