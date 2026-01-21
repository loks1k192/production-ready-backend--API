package config

import "time"

// Config holds all application configuration derived from environment variables.
type Config struct {
	Env  string
	HTTP HTTPConfig
	DB   DBConfig
	Auth AuthConfig
	Log  LogConfig
}

// HTTPConfig describes HTTP server settings.
type HTTPConfig struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DBConfig describes database settings.
type DBConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	PingTimeout     time.Duration
}

// AuthConfig describes authentication settings.
type AuthConfig struct {
	JWTSecret string
	TokenTTL  time.Duration
}

// LogConfig describes logger settings.
type LogConfig struct {
	Level string
}
