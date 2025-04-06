package router

import (
	"leaderboard-service/handlers"
	"leaderboard-service/middleware"

	"github.com/go-chi/chi/v5"
)

func init() {
	// Register protected routes
	RegisterProtectedRoutes(setupParticipantRoutes)
}

// setupParticipantRoutes configures all routes related to participants
func setupParticipantRoutes(r chi.Router) {
	// Participant routes
	r.Route("/participants", func(r chi.Router) {
		// Public participant endpoints - any authenticated user can access
		r.Get("/", handlers.ListParticipants)
		r.Get("/{id}", handlers.GetParticipant)

		// Nested routes for participant's metric values
		r.Get("/{id}/metric-values", handlers.ListMetricValues) // Get all metric values for a specific participant

		// Admin-only participant endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
			r.Post("/", handlers.CreateParticipant)
			r.Put("/{id}", handlers.UpdateParticipant)
			r.Delete("/{id}", handlers.DeleteParticipant)

			// Admin-only nested routes
			r.Post("/{id}/metric-values", handlers.CreateMetricValue) // Record a new metric value for a participant
		})
	})
}
