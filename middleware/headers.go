// middleware/headers.go
package middleware

import (
	"net/http"
)

// SecureHeadersMiddleware adds security headers to all responses
func SecureHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Block MIME-type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent XSS in older browsers
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Control referrer leakage
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Enable browser HSTS (force HTTPS) — 1 year
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy (CSP) — allow only same origin
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'")

		// Pass to next handler
		next.ServeHTTP(w, r)
	})
}
