package service

import (
	"context"
	"errors"
	"testing"

	"github.com/loks1k192/go-backend/internal/models"
	"github.com/loks1k192/go-backend/internal/repository"
)

type mockUserRepo struct {
	create     func(ctx context.Context, user models.User) (models.User, error)
	getByID    func(ctx context.Context, id int64) (models.User, error)
	getByEmail func(ctx context.Context, email string) (models.User, error)
	update     func(ctx context.Context, user models.User) (models.User, error)
	delete     func(ctx context.Context, id int64) error
}

func (m *mockUserRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	return m.create(ctx, user)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (models.User, error) {
	return m.getByID(ctx, id)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	return m.getByEmail(ctx, email)
}

func (m *mockUserRepo) Update(ctx context.Context, user models.User) (models.User, error) {
	return m.update(ctx, user)
}

func (m *mockUserRepo) Delete(ctx context.Context, id int64) error {
	return m.delete(ctx, id)
}

type fakeHasher struct{}

func (f *fakeHasher) Hash(_ context.Context, password string) (string, error) {
	return "hash:" + password, nil
}

func (f *fakeHasher) Compare(_ context.Context, hash, password string) error {
	if hash != "hash:"+password {
		return errors.New("mismatch")
	}
	return nil
}

func TestUserServiceCreate(t *testing.T) {
	repo := &mockUserRepo{
		create: func(ctx context.Context, user models.User) (models.User, error) {
			if user.Email != "user@example.com" {
				t.Fatalf("unexpected email: %s", user.Email)
			}
			if user.HashedPassword != "hash:secret123" {
				t.Fatalf("unexpected hash: %s", user.HashedPassword)
			}
			if user.Name != "Test User" {
				t.Fatalf("unexpected name: %s", user.Name)
			}
			return models.User{ID: 1, Email: user.Email, Name: user.Name}, nil
		},
	}

	service := NewUserService(repo, &fakeHasher{})
	user, err := service.Create(context.Background(), CreateInput{
		Email:    "USER@example.com",
		Password: "secret123",
		Name:     "Test User",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != 1 {
		t.Fatalf("expected id 1, got %d", user.ID)
	}
}

func TestUserServiceCreateInvalidEmail(t *testing.T) {
	service := NewUserService(&mockUserRepo{}, &fakeHasher{})

	_, err := service.Create(context.Background(), CreateInput{
		Email:    "invalid",
		Password: "secret123",
		Name:     "Test User",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
}

func TestUserServiceGetByIDNotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByID: func(ctx context.Context, id int64) (models.User, error) {
			return models.User{}, repository.ErrNotFound
		},
	}

	service := NewUserService(repo, &fakeHasher{})
	_, err := service.GetByID(context.Background(), 42)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestUserServiceAuthenticate(t *testing.T) {
	repo := &mockUserRepo{
		getByEmail: func(ctx context.Context, email string) (models.User, error) {
			if email != "user@example.com" {
				t.Fatalf("unexpected email: %s", email)
			}
			return models.User{
				ID:             10,
				Email:          email,
				HashedPassword: "hash:secret123",
			}, nil
		},
	}

	service := NewUserService(repo, &fakeHasher{})
	user, err := service.Authenticate(context.Background(), "user@example.com", "secret123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != 10 {
		t.Fatalf("expected id 10, got %d", user.ID)
	}
}
