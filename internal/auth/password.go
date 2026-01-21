package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher implements password hashing using bcrypt.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher returns a BcryptHasher with a default cost.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	return &BcryptHasher{cost: cost}
}

// Hash generates a bcrypt hash for a password.
func (h *BcryptHasher) Hash(_ context.Context, password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// Compare checks a password against a bcrypt hash.
func (h *BcryptHasher) Compare(_ context.Context, hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
