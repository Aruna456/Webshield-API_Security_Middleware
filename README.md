# WebShield – Secure Go HTTP Middleware

**A lightweight, pluggable security middleware for Go web APIs**  
**Protects against OWASP Top 10** with **zero external dependencies**.

---

## Features

| Feature | OWASP Protection | Status |
|-------|------------------|--------|
| JWT Authentication | A02:2021 | Done |
| Rate Limiting | A06:2021 | Done |
| Input Sanitization | A03:2021 | Done |
| Secure Headers | A05:2021 | Done |
| JSON Audit Logging | A09:2021 | Done |


---
## Usage 

```go
import "github.com/Aruna456/webshield/middleware"
router.Handle("/api/users", middleware.JWTMiddleware(http.HandlerFunc(handler)))
```

---
## Quick Start

```bash
# 1. Clone & enter
git clone https://github.com/yourusername/webshield.git
cd webshield

# 2. Run
go run main.go
```
## Middleware chain

```go
handler := middleware.JSONLoggingMiddleware(
    middleware.SecureHeadersMiddleware(
        middleware.RateLimitMiddleware(
            middleware.JWTMiddleware(
                middleware.SanitizeMiddleware(
                    http.HandlerFunc(usersHandler),
                    middleware.WithQuery(),
                    middleware.WithBody(),
                    middleware.WithAllowedFields(map[string][]string{
                        "name":  {"string"},
                        "age":   {"numeric"},
                        "email": {"string"},
                    }),
                ),
            ),
            10, time.Minute,
        ),
    ),
)
``` 
## Test with postman

```bash
# 1.Valid request

GET /api/users
Authorization: Bearer <jwt-with-sub>

# Result → 200 OK + JSON log with user_id

# 2. Invalid JWT

# Result → 401 Unauthorized

# 3. Rate Limit

# Result → Send 11 requests → 429 Too Many Requests

# 4. SQL/XSS Injection

{ "name": "aruna'); DROP TABLE users; --" } (or) { "name": "<script> Hackher </script>" }

# Result → 422 + debug log: &#39;); DROP... (or) &#39Hacker&#39
```
---
## Security Headers (Verified)

```http
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
Content-Security-Policy: default-src 'self'; ...
Strict-Transport-Security: max-age=31536000; ...
```
---

## Sample JSON Audit Log

```json
[AUDIT] {
  "timestamp": "2025-10-29T12:00:00Z",
  "method": "GET",
  "path": "/api/users",
  "client_ip": "127.0.0.1",
  "user_id": "alice123",
  "status": 200,
  "latency_ms": 3
}
```
---
## Tech stack

- Go 1.21+
- dgrijalva/jwt-go
- Standard Library Only

---
--- 

        Made with ❤️ @Aruna456 