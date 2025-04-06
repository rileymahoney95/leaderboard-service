package router

import (
	"leaderboard-service/handlers"
	"leaderboard-service/middleware"

	"github.com/go-chi/chi/v5"
)

func init() {
	// Register protected routes
	RegisterProtectedRoutes(setupLeaderboardRoutes)
}

// setupLeaderboardRoutes configures all routes related to leaderboards
func setupLeaderboardRoutes(r chi.Router) {
	// Leaderboard routes
	r.Route("/leaderboards", func(r chi.Router) {
		// Public leaderboard endpoints - any authenticated user can access
		r.Get("/", handlers.ListLeaderboards)
		r.Get("/{id}", handlers.GetLeaderboard)

		// Nested routes for leaderboard entries
		r.Get("/{id}/entries", handlers.ListLeaderboardEntries) // Get all entries for a specific leaderboard

		// Nested routes for leaderboard metrics
		r.Get("/{id}/metrics", handlers.ListLeaderboardMetrics) // Get all metrics for a specific leaderboard

		// Admin-only leaderboard endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
			r.Post("/", handlers.CreateLeaderboard)
			r.Put("/{id}", handlers.UpdateLeaderboard)
			r.Delete("/{id}", handlers.DeleteLeaderboard)

			// Admin-only nested routes
			r.Post("/{id}/entries", handlers.CreateLeaderboardEntry)  // Create entry for a specific leaderboard
			r.Post("/{id}/metrics", handlers.CreateLeaderboardMetric) // Associate a metric with a leaderboard
		})
	})
}
