package router

import (
	"leaderboard-service/handlers"
	"leaderboard-service/middleware"

	"github.com/go-chi/chi/v5"
)

func init() {
	// Register protected routes
	RegisterProtectedRoutes(setupMetricRoutes)
}

// setupMetricRoutes configures all routes related to metrics
func setupMetricRoutes(r chi.Router) {
	// Metric routes
	r.Route("/metrics", func(r chi.Router) {
		// Public metric endpoints - any authenticated user can access
		r.Get("/", handlers.ListMetrics)
		r.Get("/{id}", handlers.GetMetric)

		// Nested routes for metric values
		r.Get("/{id}/values", handlers.ListMetricValues) // Get all values for a specific metric

		// Admin-only metric endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
			r.Post("/", handlers.CreateMetric)
			r.Put("/{id}", handlers.UpdateMetric)
			r.Delete("/{id}", handlers.DeleteMetric)

			// Admin-only nested routes
			r.Post("/{id}/values", handlers.CreateMetricValue) // Create a new value for a specific metric
		})
	})
}
