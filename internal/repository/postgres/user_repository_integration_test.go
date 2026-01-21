//go:build integration
// +build integration

package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/loks1k192/go-backend/internal/models"
)

func TestUserRepositoryCRUD(t *testing.T) {
	ctx := context.Background()

	container, dsn := startPostgresContainer(t, ctx)
	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	if err := createUsersTable(ctx, db); err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	repo := NewUserRepository(db)
	user, err := repo.Create(ctx, models.User{
		Email:          "user@example.com",
		HashedPassword: "hashed",
		Name:           "User",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	loaded, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if loaded.Email != user.Email {
		t.Fatalf("expected email %s, got %s", user.Email, loaded.Email)
	}

	updated, err := repo.Update(ctx, models.User{
		ID:             user.ID,
		Email:          "updated@example.com",
		HashedPassword: "hashed",
		Name:           "Updated",
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Email != "updated@example.com" {
		t.Fatalf("expected updated email, got %s", updated.Email)
	}

	if err := repo.Delete(ctx, user.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}

func startPostgresContainer(t *testing.T, ctx context.Context) (testcontainers.Container, string) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       "app_test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("container start failed: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("container host failed: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("container port failed: %v", err)
	}

	dsn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/app_test?sslmode=disable", host, port.Port())
	return container, dsn
}

func createUsersTable(ctx context.Context, db *sqlx.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			hashed_password TEXT NOT NULL,
			name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`

	_, err := db.ExecContext(ctx, query)
	return err
}
