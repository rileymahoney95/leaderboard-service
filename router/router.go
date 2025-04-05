package router

import (
	"net/http"

	"leaderboard-service/handlers"
	"leaderboard-service/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router() http.Handler {
	r := chi.NewRouter()

	// Basic middleware for all routes
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.RequestLogger) // Our custom request logger
	r.Use(chimiddleware.Recoverer)

	// Public routes
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

	// Protected routes - require JWT authentication
	r.Group(func(r chi.Router) {
		// Apply JWT authentication middleware
		r.Use(middleware.JWTAuth)

		// Leaderboard routes
		r.Route("/leaderboards", func(r chi.Router) {
			// Public leaderboard endpoints - any authenticated user can access
			r.Get("/", handlers.ListLeaderboards)
			r.Get("/{id}", handlers.GetLeaderboard)

			// Admin-only leaderboard endpoints
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
				r.Post("/", handlers.CreateLeaderboard)
				r.Put("/{id}", handlers.UpdateLeaderboard)
				r.Delete("/{id}", handlers.DeleteLeaderboard)
			})
		})

		// Participant routes
		r.Route("/participants", func(r chi.Router) {
			// Public participant endpoints - any authenticated user can access
			r.Get("/", handlers.ListParticipants)
			r.Get("/{id}", handlers.GetParticipant)

			// Admin-only participant endpoints
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
				r.Post("/", handlers.CreateParticipant)
				r.Put("/{id}", handlers.UpdateParticipant)
				r.Delete("/{id}", handlers.DeleteParticipant)
			})
		})

		// Leaderboard Entry routes
		r.Route("/leaderboard-entries", func(r chi.Router) {
			// Public leaderboard entry endpoints - any authenticated user can access
			r.Get("/", handlers.ListLeaderboardEntries)
			r.Get("/{id}", handlers.GetLeaderboardEntry)

			// Admin-only leaderboard entry endpoints
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
				r.Post("/", handlers.CreateLeaderboardEntry)
				r.Put("/{id}", handlers.UpdateLeaderboardEntry)
				r.Delete("/{id}", handlers.DeleteLeaderboardEntry)
			})
		})

		// Leaderboard Metric routes
		r.Route("/leaderboard-metrics", func(r chi.Router) {
			// Public leaderboard metric endpoints - any authenticated user can access
			r.Get("/", handlers.ListLeaderboardMetrics)
			r.Get("/{id}", handlers.GetLeaderboardMetric)

			// Admin-only leaderboard metric endpoints
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
				r.Post("/", handlers.CreateLeaderboardMetric)
				r.Put("/{id}", handlers.UpdateLeaderboardMetric)
				r.Delete("/{id}", handlers.DeleteLeaderboardMetric)
			})
		})
	})

	return r
}
