// This file implements simple in-memory rate limiting per IP.
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/Aruna456/webshield/internal/config"
)

// rateLimiter struct holds visit counts per IP (in-memory for simplicity).
type rateLimiter struct {
	visits map[string]*visit
	mu     sync.Mutex
}

type visit struct {
	count     int
	lastReset time.Time
}

// Global limiter instance.
var limiter = &rateLimiter{visits: make(map[string]*visit)}

// RateLimit returns a middleware that limits requests per IP.
// If exceeded, returns 429 Too Many Requests.
func RateLimit(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr // Use IP for limiting; in prod, consider X-Forwarded-For.

			limiter.mu.Lock()
			v, ok := limiter.visits[ip]
			if !ok {
				v = &visit{}
				limiter.visits[ip] = v
			}

			now := time.Now()
			if now.Sub(v.lastReset) > time.Duration(cfg.RateLimitDuration)*time.Second {
				v.count = 0
				v.lastReset = now
			}

			if v.count >= cfg.RateLimitRequests {
				limiter.mu.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			v.count++
			limiter.mu.Unlock()

			// Pass to next.
			next.ServeHTTP(w, r)
		})
	}
}
