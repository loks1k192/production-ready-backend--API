package httpapi

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type contextKey string

const authContextKey contextKey = "auth"

// AuthInfo contains data extracted from a JWT token.
type AuthInfo struct {
	UserID int64
	Email  string
}

// AuthMiddleware validates JWT tokens and injects auth info into context.
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing token")
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := h.tokens.Parse(tokenString)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			h.logger.Warn("invalid token subject", zap.Error(err))
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), authContextKey, AuthInfo{
			UserID: userID,
			Email:  claims.Email,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthFromContext retrieves auth info from context.
func AuthFromContext(ctx context.Context) (AuthInfo, bool) {
	value := ctx.Value(authContextKey)
	if value == nil {
		return AuthInfo{}, false
	}
	info, ok := value.(AuthInfo)
	return info, ok
}
