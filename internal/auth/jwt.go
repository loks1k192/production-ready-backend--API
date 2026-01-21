package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/loks1k192/go-backend/internal/models"
)

// TokenManager issues and validates JWT tokens.
type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

// Claims describes the JWT claims payload.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// NewTokenManager creates a new JWT manager.
func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

// Generate creates a signed JWT for a user.
func (m *TokenManager) Generate(user models.User) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Parse validates a token and returns the claims.
func (m *TokenManager) Parse(tokenString string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return claims, nil
}
