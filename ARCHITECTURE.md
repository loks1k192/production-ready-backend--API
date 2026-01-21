# Architecture

This document captures the initial architecture, module boundaries, and a
step-by-step implementation plan for the Go backend.

## Goals
- Production-ready HTTP/JSON API with CRUD + health + metrics.
- Clean modular architecture (handlers, services, repositories, models).
- Postgres storage via `sqlx`.
- JWT auth for protected endpoints.
- Tests: unit + integration (Postgres).
- CI for formatting, linting, vet, tests, and Docker build.
- Docker + docker-compose + minimal K8s manifests.
- OpenAPI spec and docs.

## High-level layout
- `cmd/api/`: entrypoint, HTTP server bootstrap.
- `internal/config/`: environment-based config.
- `internal/db/`: DB connection, migrations, repo base.
- `internal/models/`: domain models.
- `internal/repository/`: SQL queries and persistence.
- `internal/service/`: business logic, validation, error mapping.
- `internal/http/`: handlers, routing, middleware.
- `internal/auth/`: JWT auth and password hashing.
- `internal/metrics/`: Prometheus setup.
- `pkg/`: reusable public packages (minimal).
- `configs/`: example env files.
- `migrations/`: SQL migration files.
- `docs/`: OpenAPI and documentation assets.

## Request flow
1. HTTP request -> router -> middleware (logging, auth, metrics, timeouts).
2. Handler validates input and calls service with `context.Context`.
3. Service applies business rules and calls repository.
4. Repository uses `sqlx` and returns domain models.
5. Handler maps result to JSON and HTTP status.

## Error handling
- Typed errors in service/repo layer.
- HTTP layer maps errors to status codes and JSON error shape.
- Errors wrapped with context for debugging.

## Testing strategy
- Unit tests for services with mocked repositories.
- Integration tests against Postgres via docker-compose or testcontainers.
- Health/metrics endpoints smoke tests.

## Implementation plan (commits)
1. Add this architecture plan.
2. Init Go module and repo metadata.
3. Add project skeleton and configs.
4. Implement DB layer and migrations.
5. Implement user model + repository.
6. Implement services and validation.
7. Implement handlers + routing + health + metrics.
8. Add JWT auth and middleware.
9. Add unit + integration tests.
10. Add CI, Docker, K8s.
11. Add OpenAPI + README + docs.
