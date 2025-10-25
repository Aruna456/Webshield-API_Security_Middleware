// This file implements JWT validation middleware.
package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Debug]Request received: %s %s", r.Method, r.URL.Path)
		//Get Auhorization token
		authHeader := r.Header.Get("Authorization")
		log.Printf("[Debug]Authorization header: %s", authHeader)
		if authHeader == "" {
			log.Println("[Error]Missing Authorization header")
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		//check for bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("[Error]Invalid Authorization header: Bearer token required")
			http.Error(w, "Invalid Authorization header: Bearer token required", http.StatusUnauthorized)
			return
		}

		// Extract Token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("[Debug]Token: %s", tokenStr)

		// Parse and Validate token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

			// Validate signing Method (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("[Error]Invalid signing method: %v", token.Method)
				return nil, http.ErrNotSupported
			}
			// Secret Key (in production, use env vars or secret management)
			return []byte("a-string-secret-at-least-256-bits-long"), nil
		})

		if err != nil || !token.Valid {
			log.Printf("[Error]Token validation failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		log.Println("[Debug]Token validated successfully")
		// Token is valid, pass to next hander
		next.ServeHTTP(w, r)
	})
}
