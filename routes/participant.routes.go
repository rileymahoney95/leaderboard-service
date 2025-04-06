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
	participantHandler := handlers.NewParticipantHandler()
	metricValueHandler := handlers.NewMetricValueHandler()

	// Participant routes
	r.Route("/participants", func(r chi.Router) {
		// Public participant endpoints - any authenticated user can access
		r.Get("/", participantHandler.ListParticipants)
		r.Get("/{id}", participantHandler.GetParticipant)

		// Nested routes for participant's metric values
		r.Get("/{id}/metric-values", metricValueHandler.ListMetricValues) // Get all metric values for a specific participant

		// Admin-only participant endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAnyRole(middleware.RoleAdmin, middleware.RoleModerator))
			r.Post("/", participantHandler.CreateParticipant)
			r.Put("/{id}", participantHandler.UpdateParticipant)
			r.Delete("/{id}", participantHandler.DeleteParticipant)

			// Admin-only nested routes
			r.Post("/{id}/metric-values", metricValueHandler.CreateMetricValue) // Record a new metric value for a participant
		})
	})
}
