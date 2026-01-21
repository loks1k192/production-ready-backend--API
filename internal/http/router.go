package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter registers routes and middleware.
func NewRouter(handler *Handler) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Get("/healthz", handler.Health)
	router.Post("/auth/login", handler.Login)

	router.Route("/users", func(r chi.Router) {
		r.Use(handler.AuthMiddleware)
		r.Post("/", handler.CreateUser)
		r.Get("/{id}", handler.GetUser)
		r.Put("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
	})

	return router
}
