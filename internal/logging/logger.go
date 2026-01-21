package logging

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/loks1k192/go-backend/internal/config"
)

// New creates a zap logger based on environment and level settings.
func New(env string, cfg config.LogConfig) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if parsed, err := zapcore.ParseLevel(strings.ToLower(cfg.Level)); err == nil {
		level = parsed
	}

	var zapCfg zap.Config
	if env == "prod" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	return zapCfg.Build()
}
