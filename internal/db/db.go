package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/loks1k192/go-backend/internal/config"
)

// New opens a Postgres connection, configures pooling, and validates connectivity.
func New(ctx context.Context, cfg config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		//nolint:errcheck
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

// Close tries to close the DB with a timeout to avoid hanging shutdowns.
func Close(ctx context.Context, db *sqlx.DB, timeout time.Duration) error {
	closeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- db.Close()
	}()

	select {
	case err := <-done:
		return err
	case <-closeCtx.Done():
		return closeCtx.Err()
	}
}
