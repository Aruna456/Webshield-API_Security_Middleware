// This file implements JWT validation middleware.
package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/Aruna456/webshield/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// JWT returns a middleware function that validates JWT from the Authorization header.
// It checks signature and expiry; returns 401 if invalid.
func JWT(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "Missing or invalid JWT", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			// Parse and validate the JWT.
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid JWT", http.StatusUnauthorized)
				return
			}

			// Check expiry.
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if exp, ok := claims["exp"].(float64); ok {
					if time.Unix(int64(exp), 0).Before(time.Now()) {
						http.Error(w, "JWT expired", http.StatusUnauthorized)
						return
					}
				}
			}

			// If valid, pass to next.
			next.ServeHTTP(w, r)
		})
	}
}
