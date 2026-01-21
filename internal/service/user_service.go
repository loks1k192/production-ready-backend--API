package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/loks1k192/go-backend/internal/models"
	"github.com/loks1k192/go-backend/internal/repository"
)

// UserService provides business logic for users.
type UserService struct {
	repo   repository.UserRepository
	hasher PasswordHasher
}

// NewUserService creates a user service with required dependencies.
func NewUserService(repo repository.UserRepository, hasher PasswordHasher) *UserService {
	return &UserService{repo: repo, hasher: hasher}
}

// CreateInput captures user creation data.
type CreateInput struct {
	Email    string
	Password string
	Name     string
}

// UpdateInput captures user update data.
type UpdateInput struct {
	ID       int64
	Email    string
	Password string
	Name     string
}

// Create creates a new user after validation and password hashing.
func (s *UserService) Create(ctx context.Context, input CreateInput) (models.User, error) {
	if err := validateEmail(input.Email); err != nil {
		return models.User{}, err
	}
	if err := validateName(input.Name); err != nil {
		return models.User{}, err
	}
	if err := validatePassword(input.Password); err != nil {
		return models.User{}, err
	}

	hashed, err := s.hasher.Hash(ctx, input.Password)
	if err != nil {
		return models.User{}, err
	}

	return s.repo.Create(ctx, models.User{
		Email:          strings.ToLower(strings.TrimSpace(input.Email)),
		HashedPassword: hashed,
		Name:           strings.TrimSpace(input.Name),
	})
}

// GetByID returns a user by ID.
func (s *UserService) GetByID(ctx context.Context, id int64) (models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

// GetByEmail returns a user by email.
func (s *UserService) GetByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := s.repo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

// Authenticate validates user credentials and returns the user on success.
func (s *UserService) Authenticate(ctx context.Context, email, password string) (models.User, error) {
	if err := validateEmail(email); err != nil {
		return models.User{}, ErrInvalidInput
	}
	if strings.TrimSpace(password) == "" {
		return models.User{}, ErrInvalidInput
	}

	user, err := s.repo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.User{}, ErrInvalidInput
		}
		return models.User{}, err
	}

	if err := s.hasher.Compare(ctx, user.HashedPassword, password); err != nil {
		return models.User{}, ErrInvalidInput
	}

	return user, nil
}

// Update updates user fields after validation and password hashing.
func (s *UserService) Update(ctx context.Context, input UpdateInput) (models.User, error) {
	if input.ID <= 0 {
		return models.User{}, ErrInvalidInput
	}
	if err := validateEmail(input.Email); err != nil {
		return models.User{}, err
	}
	if err := validateName(input.Name); err != nil {
		return models.User{}, err
	}
	if err := validatePassword(input.Password); err != nil {
		return models.User{}, err
	}

	hashed, err := s.hasher.Hash(ctx, input.Password)
	if err != nil {
		return models.User{}, err
	}

	user, err := s.repo.Update(ctx, models.User{
		ID:             input.ID,
		Email:          strings.ToLower(strings.TrimSpace(input.Email)),
		HashedPassword: hashed,
		Name:           strings.TrimSpace(input.Name),
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

// Delete removes a user by ID.
func (s *UserService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidInput
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func validateEmail(email string) error {
	trimmed := strings.TrimSpace(email)
	if trimmed == "" {
		return ErrInvalidInput
	}
	if _, err := mail.ParseAddress(trimmed); err != nil {
		return ErrInvalidInput
	}

	return nil
}

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidInput
	}

	return nil
}

func validatePassword(password string) error {
	if strings.TrimSpace(password) == "" || len(password) < 8 {
		return ErrInvalidInput
	}

	return nil
}
