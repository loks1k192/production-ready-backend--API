package repository

import (
	"context"

	"github.com/loks1k192/go-backend/internal/models"
)

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	GetByID(ctx context.Context, id int64) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	Delete(ctx context.Context, id int64) error
}
