package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Aruna456/webshield/middleware"
)

// Sample handler for testing
func userHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Users API!")
}

func main() {
	// Creating a  Router
	mux := http.NewServeMux()
	// Chain middleware: Rate Limiting -> JWT Authentication -> Handler
	// mux.Handle("/api/users", middleware.JWTMiddleware(http.HandlerFunc(userHandler)))
	handler := middleware.RateLimitMiddleware(
		middleware.JWTMiddleware(http.HandlerFunc(userHandler)),
		10,          // 10 requests
		time.Minute, // per minute
	)
	//Register endpoints
	mux.Handle("/api/users", handler)
	// start server
	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed:%v", err)
	}
}
