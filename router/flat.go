package router

import (
	"leaderboard-service/handlers"
	"leaderboard-service/middleware"

	"github.com/go-chi/chi/v5"
)

func init() {
	// Register protected routes
	RegisterProtectedRoutes(setupFlatRoutes)
}

// setupFlatRoutes configures all routes for backward compatibility
func setupFlatRoutes(r chi.Router) {
	// Leaderboard Entry routes (flat, original)
	r.Route("/leaderboard-entries", func(r chi.Router) {
		// Public leaderboard entry endpoints
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

	// Leaderboard Metric routes (flat, original)
	r.Route("/leaderboard-metrics", func(r chi.Router) {
		// Public leaderboard metric endpoints
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

	// Metric Value routes (flat, original)
	r.Route("/metric-values", func(r chi.Router) {
		// Public metric value endpoints
		r.Get("/", handlers.ListMetricValues)
		r.Get("/{id}", handlers.GetMetricValue)

		// Admin-only metric value endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
			r.Post("/", handlers.CreateMetricValue)
			r.Put("/{id}", handlers.UpdateMetricValue)
			r.Delete("/{id}", handlers.DeleteMetricValue)
		})
	})
}
