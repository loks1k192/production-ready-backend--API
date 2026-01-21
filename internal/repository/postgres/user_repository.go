package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/loks1k192/go-backend/internal/models"
	"github.com/loks1k192/go-backend/internal/repository"
)

// UserRepository implements repository.UserRepository with Postgres.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository returns a Postgres-backed user repository.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create persists a new user and returns the created record.
func (r *UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	const query = `
		INSERT INTO users (email, hashed_password, name)
		VALUES ($1, $2, $3)
		RETURNING id, email, hashed_password, name, created_at;
	`

	var created models.User
	if err := r.db.GetContext(ctx, &created, query, user.Email, user.HashedPassword, user.Name); err != nil {
		return models.User{}, err
	}

	return created, nil
}

// GetByID loads a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (models.User, error) {
	const query = `
		SELECT id, email, hashed_password, name, created_at
		FROM users
		WHERE id = $1;
	`

	var user models.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repository.ErrNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

// GetByEmail loads a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	const query = `
		SELECT id, email, hashed_password, name, created_at
		FROM users
		WHERE email = $1;
	`

	var user models.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repository.ErrNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

// Update modifies user fields and returns the updated record.
func (r *UserRepository) Update(ctx context.Context, user models.User) (models.User, error) {
	const query = `
		UPDATE users
		SET email = $1,
			hashed_password = $2,
			name = $3
		WHERE id = $4
		RETURNING id, email, hashed_password, name, created_at;
	`

	var updated models.User
	if err := r.db.GetContext(ctx, &updated, query, user.Email, user.HashedPassword, user.Name, user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repository.ErrNotFound
		}
		return models.User{}, err
	}

	return updated, nil
}

// Delete removes a user by ID.
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM users
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
