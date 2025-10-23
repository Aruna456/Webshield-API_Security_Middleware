// internal/middleware/logging.go
// Request logging middleware with correlation ID (UUID v4, standard library only)

package middleware

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Aruna456/webshield/internal/config"
)

// Logging returns a middleware that logs every request with:
// - Correlation ID (UUID)
// - Method, Path, Status, Duration
func Logging(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate correlation ID if not present
			corrID := r.Header.Get(cfg.CorrelationIDHeader)
			if corrID == "" {
				corrID = generateUUID()
			}
			r.Header.Set(cfg.CorrelationIDHeader, corrID)

			// Start timer
			start := time.Now()

			// Wrap response writer to capture status code
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call next handler
			next.ServeHTTP(ww, r)

			// Log after response
			duration := time.Since(start)
			log.Printf("[%s] %s %s %d %v", corrID, r.Method, r.URL.Path, ww.statusCode, duration)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// generateUUID creates a simple UUID v4 using crypto/rand (standard library)
func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "unknown-id"
	}
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant RFC 4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
