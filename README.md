# Go Backend API

[![CI](https://github.com/loks1k192/production-ready-backend--API/actions/workflows/ci.yml/badge.svg)](https://github.com/loks1k192/production-ready-backend--API/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/Go-1.24-blue)

Production-ready Go backend with REST API, JWT auth, PostgreSQL, metrics, and CI.

## Features
- REST CRUD for users
- JWT auth (login endpoint + middleware)
- PostgreSQL storage (sqlx)
- Prometheus metrics at `/metrics`
- Health check at `/healthz`
- Structured logging (zap)
- Docker, docker-compose, and Kubernetes manifests
- Unit + integration tests (Postgres)

## Prerequisites
- Go 1.24+
- Docker + Docker Compose (for local stack and integration tests)
- PostgreSQL 16+ (if running locally without Docker)
- `golang-migrate` for database migrations

## Quickstart (Docker Compose)
```bash
docker-compose up --build
```

Run migrations from your host:
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/app?sslmode=disable" up
```

## Development (local, no Docker)
```bash
export APP_ENV=dev
export HTTP_ADDR=:8080
export DB_DSN="postgres://postgres:postgres@localhost:5432/app?sslmode=disable"
export AUTH_JWT_SECRET="change-me"
export AUTH_TOKEN_TTL=1h
export LOG_LEVEL=info

migrate -path migrations -database "$DB_DSN" up
go run ./cmd/api
```

## Configuration
Use environment variables (see `configs/env.example`).

Key variables:
- `APP_ENV`: `dev` or `prod`
- `HTTP_ADDR`: `:8080`
- `DB_DSN`: `postgres://...`
- `AUTH_JWT_SECRET`: JWT secret key
- `AUTH_TOKEN_TTL`: token TTL (e.g. `1h`)
- `LOG_LEVEL`: `debug|info|warn|error`

## API quick start
### Login
```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
```

### Create user
```bash
curl -s -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123","name":"User"}'
```

### Get user
```bash
curl -s http://localhost:8080/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

## OpenAPI
OpenAPI spec is located at `docs/openapi.yaml`.

Swagger UI (Docker):
```bash
docker run --rm -p 8081:8080 \
  -e SWAGGER_JSON=/openapi.yaml \
  -v "$PWD/docs/openapi.yaml:/openapi.yaml" \
  swaggerapi/swagger-ui
```

## Migrations
Migrations are in `migrations/`. You can run them with `golang-migrate`:
```bash
migrate -path migrations -database "$DB_DSN" up
```

## Testing & Linting
```bash
gofmt -w .
go vet ./...
go test ./...
go test -tags=integration ./...
```

Install and run golangci-lint (same as CI):
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4
golangci-lint run
```

## CI
Workflow: `.github/workflows/ci.yml`
- gofmt, go vet, unit tests, integration tests
- golangci-lint (built with the Go toolchain from `go.mod`)
- Docker build

## Troubleshooting
- `golangci-lint` version mismatch: install it via `go install` using the same Go version as `go.mod`, or rely on CI which builds it with the project toolchain.

## Contributing
Issues and PRs are welcome. Please keep code formatted with `gofmt` and ensure tests and lint pass.

## License
Not specified yet.

## Contacts
Maintainer: [loks1k192](https://github.com/loks1k192)
