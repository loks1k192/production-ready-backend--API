# Go Backend API

Production-ready Go backend with REST API, JWT auth, PostgreSQL, metrics, and CI.

## Features
- REST CRUD for users
- JWT auth (login endpoint + middleware)
- PostgreSQL storage (sqlx)
- Prometheus metrics at `/metrics`
- Health check at `/healthz`
- Structured logging (zap)
- Docker, docker-compose, and Kubernetes manifests
- Unit + integration tests

## Local development
### Docker Compose
```bash
docker-compose up --build
```

### Go run
```bash
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
  -H "Authorization: Bearer $TOKEN" \
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

## Testing
```bash
go test ./...
go test -tags=integration ./...
```
