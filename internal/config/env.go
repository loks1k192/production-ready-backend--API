package config

import (
	"os"
	"strconv"
	"time"
)

// Load reads configuration from environment variables with sensible defaults.
func Load() (Config, error) {
	return Config{
		Env: getString("APP_ENV", "dev"),
		HTTP: HTTPConfig{
			Addr:            getString("HTTP_ADDR", ":8080"),
			ReadTimeout:     getDuration("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     getDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		DB: DBConfig{
			DSN:             getString("DB_DSN", "postgres://postgres:postgres@localhost:5432/app?sslmode=disable"),
			MaxOpenConns:    getInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
			PingTimeout:     getDuration("DB_PING_TIMEOUT", 5*time.Second),
		},
		Auth: AuthConfig{
			JWTSecret: getString("AUTH_JWT_SECRET", "change-me"),
			TokenTTL:  getDuration("AUTH_TOKEN_TTL", 1*time.Hour),
		},
		Log: LogConfig{
			Level: getString("LOG_LEVEL", "info"),
		},
	}, nil
}

func getString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}

	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}

	return fallback
}
