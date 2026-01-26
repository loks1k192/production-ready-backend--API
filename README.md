# Go Backend API

[![CI](https://github.com/loks1k192/production-ready-backend--API/actions/workflows/ci.yml/badge.svg)](https://github.com/loks1k192/production-ready-backend--API/actions/workflows/ci.yml)

Production-ready backend на Go с REST API, JWT-аутентификацией, PostgreSQL, метриками и CI.

Проект демонстрирует типичную архитектуру серверного приложения на Go, готового к использованию в продакшене.

## Возможности
- REST CRUD API для пользователей
- JWT-аутентификация (endpoint логина + middleware)
- Хранение данных в PostgreSQL (через `sqlx`)
- Метрики Prometheus по адресу `/metrics`
- Health check по адресу `/healthz`
- Структурированное логирование (zap)
- Docker, docker-compose и Kubernetes-манифесты
- Unit- и integration-тесты (с PostgreSQL)

## Требования
- Go версии 1.24 или выше
- Docker и Docker Compose (для локального запуска и интеграционных тестов)
- PostgreSQL 16+ (если запускать без Docker)
- `golang-migrate` для управления миграциями БД

## Быстрый старт (Docker Compose)
Запуск всего стека:
```bash
docker-compose up --build
````

Запуск миграций с хоста:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/app?sslmode=disable" up
```

## Разработка локально (без Docker)

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

## Конфигурация

Приложение настраивается через переменные окружения
(см. пример в `configs/env.example`).

Основные параметры:

* `APP_ENV` — окружение (`dev` или `prod`)
* `HTTP_ADDR` — HTTP-адрес сервера (например `:8080`)
* `DB_DSN` — строка подключения к PostgreSQL
* `AUTH_JWT_SECRET` — секретный ключ для JWT
* `AUTH_TOKEN_TTL` — время жизни токена (например `1h`)
* `LOG_LEVEL` — уровень логирования (`debug|info|warn|error`)

## Быстрый старт API

### Логин

```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
```

### Создание пользователя

```bash
curl -s -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123","name":"User"}'
```

### Получение пользователя

```bash
curl -s http://localhost:8080/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

## OpenAPI

OpenAPI-спецификация находится в `docs/openapi.yaml`.

Запуск Swagger UI через Docker:

```bash
docker run --rm -p 8081:8080 \
  -e SWAGGER_JSON=/openapi.yaml \
  -v "$PWD/docs/openapi.yaml:/openapi.yaml" \
  swaggerapi/swagger-ui
```

## Миграции

Все миграции лежат в директории `migrations/`.
Для запуска используется `golang-migrate`:

```bash
migrate -path migrations -database "$DB_DSN" up
```

## Тестирование и линтинг

```bash
gofmt -w .
go vet ./...
go test ./...
go test -tags=integration ./...
```

Установка и запуск golangci-lint (аналогично CI):

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4
golangci-lint run
```

## CI

Workflow находится в `.github/workflows/ci.yml` и включает:

* gofmt, go vet
* unit- и integration-тесты
* golangci-lint (сборка с использованием версии Go из `go.mod`)
* сборку Docker-образа

## Возможные проблемы

* Несовпадение версии `golangci-lint`:
  устанавливайте линтер через `go install`, используя ту же версию Go, что указана в `go.mod`,
  либо ориентируйтесь на CI — там линтер собирается с нужным toolchain.

