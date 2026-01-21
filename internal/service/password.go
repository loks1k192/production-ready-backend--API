package service

import "context"

// PasswordHasher abstracts password hashing.
type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, hash, password string) error
}
