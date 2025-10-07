package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Aruna456/webShield/handlers"
	"github.com/Aruna456/webShield/middleware"
)

func main() {

	// http.HandleFunc("/", handlers.HelloHandler)
	rl := middleware.NewRateLimiter(5, time.Minute)
	http.Handle("/", rl.Middleware(http.HandlerFunc(handlers.HelloHandler)))
	fmt.Println("Server Starting on port:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
