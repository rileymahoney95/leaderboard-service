package router

import (
	"net/http"

	"leaderboard-service/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// Define types for route setup functions for better organization
type RouteSetupFunc func(r chi.Router)

// RouteGroups holds all the route setup functions
type RouteGroups struct {
	Public    []RouteSetupFunc
	Protected []RouteSetupFunc
}

// Global route groups that will be populated by init() functions in other files
var routes = RouteGroups{
	Public:    []RouteSetupFunc{},
	Protected: []RouteSetupFunc{},
}

// RegisterPublicRoutes adds route setup functions to the public routes list
func RegisterPublicRoutes(setupFunc RouteSetupFunc) {
	routes.Public = append(routes.Public, setupFunc)
}

// RegisterProtectedRoutes adds route setup functions to the protected routes list
func RegisterProtectedRoutes(setupFunc RouteSetupFunc) {
	routes.Protected = append(routes.Protected, setupFunc)
}

// Router creates and configures the main router
func Router() http.Handler {
	r := chi.NewRouter()

	// Basic middleware for all routes
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.RequestLogger) // Our custom request logger
	r.Use(chimiddleware.Recoverer)

	// Mount public routes
	for _, setupFunc := range routes.Public {
		setupFunc(r)
	}

	// Protected routes - require JWT authentication
	r.Group(func(r chi.Router) {
		// Apply JWT authentication middleware
		r.Use(middleware.JWTAuth)

		// Mount all protected routes
		for _, setupFunc := range routes.Protected {
			setupFunc(r)
		}
	})

	return r
}
