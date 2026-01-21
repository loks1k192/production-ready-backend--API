package models

import "time"

// User represents the persisted user entity.
type User struct {
	ID             int64     `db:"id" json:"id"`
	Email          string    `db:"email" json:"email"`
	HashedPassword string    `db:"hashed_password" json:"-"`
	Name           string    `db:"name" json:"name"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
