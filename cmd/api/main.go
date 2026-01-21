package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/loks1k192/go-backend/internal/auth"
	"github.com/loks1k192/go-backend/internal/config"
	"github.com/loks1k192/go-backend/internal/db"
	httpapi "github.com/loks1k192/go-backend/internal/http"
	"github.com/loks1k192/go-backend/internal/logging"
	"github.com/loks1k192/go-backend/internal/repository/postgres"
	"github.com/loks1k192/go-backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger, err := logging.New(cfg.Env, cfg.Log)
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	ctx := context.Background()
	database, err := db.New(ctx, cfg.DB)
	if err != nil {
		logger.Fatal("db connection failed", zap.Error(err))
	}

	userRepo := postgres.NewUserRepository(database)
	hasher := auth.NewBcryptHasher(0)
	userService := service.NewUserService(userRepo, hasher)
	tokens := auth.NewTokenManager(cfg.Auth.JWTSecret, cfg.Auth.TokenTTL)

	handler := httpapi.NewHandler(userService, tokens, logger)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		logger.Info("http server started", zap.String("addr", cfg.HTTP.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("http shutdown failed", zap.Error(err))
	}
	if err := db.Close(shutdownCtx, database, cfg.HTTP.ShutdownTimeout); err != nil {
		logger.Error("db shutdown failed", zap.Error(err))
	}

	logger.Info("shutdown complete")
}
