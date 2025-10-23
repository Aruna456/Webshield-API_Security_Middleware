// Package config handles loading and parsing the application configuration from YAML.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config struct holds all configurable settings for the application.
// These can be toggled or adjusted in config.yaml.
type Config struct {
	Port                int    `yaml:"port"`                  // Port to run the server on.
	APIKey              string `yaml:"api_key"`               // Secret API key for authentication.
	JWTSecret           string `yaml:"jwt_secret"`            // Secret for signing/verifying JWTs.
	EnableRateLimit     bool   `yaml:"enable_rate_limit"`     // Enable/disable rate limiting.
	RateLimitRequests   int    `yaml:"rate_limit_requests"`   // Max requests per IP per duration.
	RateLimitDuration   int    `yaml:"rate_limit_duration"`   // Time window for rate limiting in seconds.
	EnableAPIKey        bool   `yaml:"enable_api_key"`        // Enable/disable API key check.
	EnableJWT           bool   `yaml:"enable_jwt"`            // Enable/disable JWT validation.
	EnableLogging       bool   `yaml:"enable_logging"`        // Enable/disable request logging.
	CorrelationIDHeader string `yaml:"correlation_id_header"` // Header name for correlation ID.
}

// Load reads the YAML file and unmarshals it into the Config struct.
// It also sets default values if not provided in the file.
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// Set defaults.
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.RateLimitRequests == 0 {
		cfg.RateLimitRequests = 60
	}
	if cfg.RateLimitDuration == 0 {
		cfg.RateLimitDuration = 60 // 1 minute.
	}
	if cfg.CorrelationIDHeader == "" {
		cfg.CorrelationIDHeader = "X-Correlation-ID"
	}
	return &cfg, nil
}
