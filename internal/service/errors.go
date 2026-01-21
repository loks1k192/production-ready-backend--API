package service

import "errors"

// ErrNotFound indicates an entity could not be located.
var ErrNotFound = errors.New("not found")

// ErrInvalidInput indicates input validation failed.
var ErrInvalidInput = errors.New("invalid input")
