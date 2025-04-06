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
	metricHandler := handlers.NewMetricHandler()
	metricValueHandler := handlers.NewMetricValueHandler()

	// Metric routes
	r.Route("/metrics", func(r chi.Router) {
		// Public metric endpoints - any authenticated user can access
		r.Get("/", metricHandler.ListMetrics)
		r.Get("/{id}", metricHandler.GetMetric)

		// Nested routes for metric values
		r.Get("/{id}/values", metricValueHandler.ListMetricValues) // Get all values for a specific metric

		// Admin-only metric endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
			r.Post("/", metricHandler.CreateMetric)
			r.Put("/{id}", metricHandler.UpdateMetric)
			r.Delete("/{id}", metricHandler.DeleteMetric)

			// Admin-only nested routes
			r.Post("/{id}/values", metricValueHandler.CreateMetricValue) // Create a new value for a specific metric
		})
	})
}
