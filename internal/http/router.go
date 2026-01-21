package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/loks1k192/go-backend/internal/metrics"
)

// NewRouter registers routes and middleware.
func NewRouter(handler *Handler) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(handler.LoggingMiddleware)

	router.Get("/healthz", handler.Health)
	router.Post("/auth/login", handler.Login)
	router.Handle("/metrics", metrics.Handler())

	router.Route("/users", func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Get("/{id}", handler.GetUser)
			r.Put("/{id}", handler.UpdateUser)
			r.Delete("/{id}", handler.DeleteUser)
		})
	})

	return router
}
