package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Aruna456/webshield/middleware"
)

// Sample handler for testing
func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Users API!")
}

func main() {
	// Creating a  Router
	mux := http.NewServeMux()
	handler := middleware.JSONLoggingMiddleware( //  LOGS EVERYTHING
		middleware.SecureHeadersMiddleware( //  HEADERS
			middleware.RateLimitMiddleware(
				middleware.JWTMiddleware(
					middleware.SanitizeMiddleware(
						http.HandlerFunc(usersHandler),
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
