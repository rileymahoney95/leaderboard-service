package router

import (
	"net/http"

	"leaderboard-service/handlers"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	// Register public routes
	RegisterPublicRoutes(setupPublicRoutes)
}

// setupPublicRoutes configures all routes that do not require authentication
func setupPublicRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		// Base routes
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Swagger documentation
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
		))

		// Authentication routes
		r.Post("/auth/login", handlers.Login)
		r.Post("/auth/register", handlers.Register)
	})
}
