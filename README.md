# WebShield: Lightweight API Security Middleware in Go

> A **beginner-friendly**, **production-ready** API security middleware with **Rate Limiting**, **API Key**, **JWT**, and **Request Logging**.  
> Built with **Go 1.25.3**, minimal dependencies, and fully tested.

---

## Features

- ✅ API Key Authentication (`X-Api-Key`)
- 🔐 JWT Validation (signature + expiry)
- 🚦 Rate Limiting per IP (in-memory, thread-safe)
- 🧾 Request Logging with Correlation ID (UUID)
- ⚙️ Configurable via `config.yaml`
- 🧩 Composable Middleware Design
- 🌐 Endpoints: `/health`, `/public`, `/protected`

---

## Project Structure


WebShield/
├── cmd/webshield/main.go
├── internal/
│ ├── config/config.go
│ ├── handlers/handlers.go
│ └── middleware/
│ ├── api_key.go
│ ├── jwt.go
│ ├── logging.go
│ ├── rate_limit.go
│ └── middleware.go
├── config.yaml
├── go.mod
├── README.md
└── tests/
└── integration_test.go



---

## Prerequisites

- **Go 1.25.3** (or newer)
- Terminal (Linux/macOS/Windows)

---

## Setup & Run

```bash
# 1. Create project
mkdir WebShield && cd WebShield

# 2. (Paste all files from the project here — structure as shown above)

# 3. Install dependencies
go mod tidy

# 4. Edit config.yaml
nano config.yaml

Example config.yaml

port: 8080
api_key: "my-super-secret-api-key-123"
jwt_secret: "my-jwt-secret-very-long-and-random"
enable_rate_limit: true
rate_limit_requests: 5
rate_limit_duration: 10
enable_api_key: true
enable_jwt: true
enable_logging: true
correlation_id_header: "X-Correlation-ID"

# 5. Run the server
go run ./cmd/webshield/main.go
```

## Limitations

- Rate limiting is in-memory (not shared between servers)

- JWT validation checks signature + expiry only

- No HTTPS (use a reverse proxy like Nginx)

- No database; secrets stored in config.yaml

## Future Improvements

🔄 Redis-based distributed rate limiting

🧩 Role-based JWT claims

📊 Prometheus metrics

🐳 Docker + Docker Compose setup

🚏 Per-route middleware with Chi route groups

---
                  Made with ❤️ for Go beginners.