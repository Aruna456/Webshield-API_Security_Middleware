package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// RateLimiter holds the state for rate limiting per client
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     int           // Tokens per interval (e.g., 10 requests)
	interval time.Duration // Time window (e.g., 1 minute)
}

// visitor tracks token bucket state for a client
type visitor struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		interval: interval,
	}
}

// RateLimitMiddleware limits requests per client IP
func RateLimitMiddleware(next http.Handler, rate int, interval time.Duration) http.Handler {
	rl := NewRateLimiter(rate, interval)

	// Background cleanup to remove old visitors
	go rl.cleanupVisitors()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Get client IP (simplified; in production, handle proxies)
		clientIP := r.RemoteAddr
		log.Printf("Rate limit check for IP: %s", clientIP)

		// Get or create visitor
		v, exists := rl.visitors[clientIP]
		if !exists {
			rl.visitors[clientIP] = &visitor{tokens: rate, lastSeen: time.Now()}
			v = rl.visitors[clientIP]
		}

		// Refill tokens based on elapsed time
		elapsed := time.Since(v.lastSeen)
		if elapsed > interval {
			v.tokens = rate
			v.lastSeen = time.Now()
		}

		// Check if tokens are available
		if v.tokens <= 0 {
			log.Printf("Rate limit exceeded for IP: %s", clientIP)
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Consume a token
		v.tokens--
		log.Printf("Tokens remaining for IP %s: %d", clientIP, v.tokens)
		v.lastSeen = time.Now()

		// Pass to next handler
		next.ServeHTTP(w, r)
	})
}

// cleanupVisitors removes stale visitor entries
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(10 * time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.interval {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
