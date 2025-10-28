// middleware/logging.go
package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// LogEntry - structured audit log
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	ClientIP  string `json:"client_ip"`
	UserID    string `json:"user_id,omitempty"`
	Status    int    `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
}

// JSONLoggingMiddleware logs every request in JSON
func JSONLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture status code
		ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(ww, r)

		// Extract user from JWT context
		userID := ""
		if uid := r.Context().Value("jwt_user_id"); uid != nil {
			userID = uid.(string)
		}

		// Build log entry
		entry := LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Method:    r.Method,
			Path:      r.URL.Path,
			ClientIP:  getIP(r),
			UserID:    userID,
			Status:    ww.status,
			LatencyMs: time.Since(start).Milliseconds(),
		}

		// Print as JSON
		jsonBytes, _ := json.Marshal(entry)
		log.Printf("[AUDIT] %s", jsonBytes)
	})
}

// responseWriter captures status
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// getIP handles X-Forwarded-For
func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}
