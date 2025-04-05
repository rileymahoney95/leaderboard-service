package router

import (
	"net/http"

	"leaderboard-service/handlers"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewRouter()

	// TODO: Add middleware

	// Base routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Leaderboard routes
	r.Route("/leaderboards", func(r chi.Router) {
		r.Post("/", handlers.CreateLeaderboard)
		r.Get("/", handlers.ListLeaderboards)
		r.Get("/{id}", handlers.GetLeaderboard)
		r.Put("/{id}", handlers.UpdateLeaderboard)
		r.Delete("/{id}", handlers.DeleteLeaderboard)
	})

	return r
}
