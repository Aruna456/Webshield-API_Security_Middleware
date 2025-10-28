package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ---------- Config & Functional Options ----------
type SanitizeConfig struct {
	ValidateQuery bool
	ValidateBody  bool
	AllowedFields map[string][]string
}

type SanitizeOption func(*SanitizeConfig)

func WithQuery() SanitizeOption {
	return func(c *SanitizeConfig) { c.ValidateQuery = true }
}
func WithBody() SanitizeOption {
	return func(c *SanitizeConfig) { c.ValidateBody = true }
}
func WithAllowedFields(m map[string][]string) SanitizeOption {
	return func(c *SanitizeConfig) { c.AllowedFields = m }
}

// ---------- Middleware ----------
func SanitizeMiddleware(next http.Handler, opts ...SanitizeOption) http.Handler {
	cfg := &SanitizeConfig{
		ValidateQuery: false,
		ValidateBody:  false,
		AllowedFields: map[string][]string{},
	}
	for _, o := range opts {
		o(cfg)
	}

	_ = validator.New() // keep import alive

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Sanitize] %s %s", r.Method, r.URL.Path)

		// ---- Query ----
		if cfg.ValidateQuery {
			if err := validateQuery(r, cfg.AllowedFields); err != nil {
				log.Printf("[Sanitize] query error: %v", err)
				http.Error(w, "Invalid query: "+err.Error(), http.StatusUnprocessableEntity)
				return
			}
			// DEBUG: print sanitized query string
			log.Printf("[DEBUG] Sanitized query: %s", r.URL.Query().Encode())
		}

		// ---- JSON body ----
		if cfg.ValidateBody && strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			cleanBody, err := validateJSONBody(r, cfg.AllowedFields)
			if err != nil {
				log.Printf("[Sanitize] body error: %v", err)
				http.Error(w, "Invalid body: "+err.Error(), http.StatusUnprocessableEntity)
				return
			}
			// DEBUG: print sanitized JSON
			prettyJSON, _ := json.MarshalIndent(cleanBody, "", "  ")
			log.Printf("[DEBUG] Sanitized JSON body:\n%s", prettyJSON)
		}

		next.ServeHTTP(w, r)
	})
}

// ---------- Query ----------
func validateQuery(r *http.Request, allowed map[string][]string) error {
	q := r.URL.Query()
	for key, vals := range q {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		rules, ok := allowed[key]
		if !ok {
			return fmt.Errorf("field %q is not allowed", key)
		}

		for _, rule := range rules {
			switch rule {
			case "numeric":
				if _, err := strconv.ParseFloat(val, 64); err != nil {
					return fmt.Errorf("field %q must be numeric", key)
				}
			}
		}
	}
	return nil
}

// ---------- JSON body ----------
func validateJSONBody(r *http.Request, allowed map[string][]string) (map[string]interface{}, error) {
	// 1. Decode
	var raw map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	defer r.Body.Close()

	// 2. Validate + clean
	clean := make(map[string]interface{})
	for k, v := range raw {
		rules, ok := allowed[k]
		if !ok {
			return nil, fmt.Errorf("field %q is not allowed", k)
		}

		for _, rule := range rules {
			switch rule {
			case "string":
				s, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("field %q must be a string", k)
				}
				s = strings.TrimSpace(s)
				s = html.EscapeString(s)
				clean[k] = s
			case "numeric":
				switch reflect.TypeOf(v).Kind() {
				case reflect.Float64, reflect.Float32, reflect.Int, reflect.Int64:
					clean[k] = v
				default:
					return nil, fmt.Errorf("field %q must be numeric", k)
				}
			}
		}
	}

	// 3. Re-inject cleaned body
	b, _ := json.Marshal(clean)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	r.ContentLength = int64(len(b))

	return clean, nil
}
