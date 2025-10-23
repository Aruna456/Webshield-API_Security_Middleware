// This file composes the middleware chain based on config.
package middleware

import (
	"net/http"

	"github.com/Aruna456/webshield/internal/config"
)

// Compose returns a slice of enabled middleware functions in order.
// Order matters: logging first, then rate limit, then auth.
func Compose(cfg *config.Config) []func(http.Handler) http.Handler {
	var mws []func(http.Handler) http.Handler

	if cfg.EnableLogging {
		mws = append(mws, Logging(cfg))
	}
	if cfg.EnableRateLimit {
		mws = append(mws, RateLimit(cfg))
	}
	if cfg.EnableAPIKey {
		mws = append(mws, APIKey(cfg))
	}
	if cfg.EnableJWT {
		mws = append(mws, JWT(cfg))
	}

	return mws
}
