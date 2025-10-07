package middleware

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimitEntry holds request count and last reset time for an IP
type RateLimitEntry struct {
	Count     int
	LastReset time.Time
}

// RateLimiter manages rate limiting for IPs
type RateLimiter struct {
	requests map[string]*RateLimitEntry
	mutex    sync.Mutex
	maxReq   int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxReq int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]*RateLimitEntry),
		maxReq:   maxReq,
		window:   window,
	}
}

// Middleware applies rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract IP (remove port if present)
		ip := r.RemoteAddr
		if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
			ip = ip[:colonIndex]
		}

		// Log the requested URL for debugging
		log.Printf("Request to %s from IP %s", r.URL.Path, ip)

		rl.mutex.Lock()
		entry, exists := rl.requests[ip]
		if !exists {
			rl.requests[ip] = &RateLimitEntry{Count: 0, LastReset: time.Now()}
			entry = rl.requests[ip]
		}

		// Check if window has expired
		if time.Since(entry.LastReset) > rl.window {
			log.Printf("Resetting count for IP %s: Window expired", ip)
			entry.Count = 0
			entry.LastReset = time.Now()
		}

		// Increment count
		entry.Count++
		log.Printf("IP %s: Count=%d, Max=%d, LastReset=%v", ip, entry.Count, rl.maxReq, entry.LastReset)

		// Check if limit exceeded
		if entry.Count > rl.maxReq {
			rl.mutex.Unlock()
			log.Printf("Rate limit exceeded for IP %s", ip)
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		rl.mutex.Unlock()
		next.ServeHTTP(w, r)
	})
}
