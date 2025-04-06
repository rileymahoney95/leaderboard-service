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

// setupFlatRoutes configures all "flat" routes (not nested under resources)
func setupFlatRoutes(r chi.Router) {
	metricValueHandler := handlers.NewMetricValueHandler()
	leaderboardEntryHandler := handlers.NewLeaderboardEntryHandler()

	// Metric Value routes (flat)
	r.Route("/metric-values", func(r chi.Router) {
		// Public endpoints
		r.Get("/", metricValueHandler.ListMetricValues)
		r.Get("/{id}", metricValueHandler.GetMetricValue)

		// Admin-only endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
			r.Post("/", metricValueHandler.CreateMetricValue)
			r.Put("/{id}", metricValueHandler.UpdateMetricValue)
			r.Delete("/{id}", metricValueHandler.DeleteMetricValue)
		})
	})

	// LeaderboardEntry routes (flat)
	r.Route("/leaderboard-entries", func(r chi.Router) {
		// Public endpoints
		r.Get("/", leaderboardEntryHandler.ListLeaderboardEntries)
		r.Get("/{id}", leaderboardEntryHandler.GetLeaderboardEntry)

		// Admin-only endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
			r.Post("/", leaderboardEntryHandler.CreateLeaderboardEntry)
			r.Put("/{id}", leaderboardEntryHandler.UpdateLeaderboardEntry)
			r.Delete("/{id}", leaderboardEntryHandler.DeleteLeaderboardEntry)
		})
	})

	// LeaderboardMetric routes (flat)
	r.Route("/leaderboard-metrics", func(r chi.Router) {
		// Public endpoints
		r.Get("/", handlers.ListLeaderboardMetrics)
		r.Get("/{id}", handlers.GetLeaderboardMetric)

		// Admin-only endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
			r.Post("/", handlers.CreateLeaderboardMetric)
			r.Put("/{id}", handlers.UpdateLeaderboardMetric)
			r.Delete("/{id}", handlers.DeleteLeaderboardMetric)
		})
	})
}
