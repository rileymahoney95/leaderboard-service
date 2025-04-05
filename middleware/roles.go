package middleware

import (
	"net/http"
)

// Role type represents user roles in the system
type Role string

const (
	// RoleAdmin represents an administrator
	RoleAdmin Role = "admin"
	// RoleModerator represents a moderator
	RoleModerator Role = "moderator"
	// RoleUser represents a regular user
	RoleUser Role = "user"
)

// RequireRole is a middleware that checks if the user has the required role
func RequireRole(requiredRole Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			claims, err := GetUserFromContext(r.Context())
			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
				return
			}

			// Check if the user has the required role
			if Role(claims.Role) != requiredRole {
				RespondWithError(w, http.StatusForbidden, "Insufficient permissions", nil)
				return
			}

			// User has the required role, proceed
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole checks if the user has any of the provided roles
func RequireAnyRole(roles ...Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			claims, err := GetUserFromContext(r.Context())
			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
				return
			}

			// Check if the user has any of the required roles
			userRole := Role(claims.Role)
			hasRequiredRole := false
			for _, role := range roles {
				if userRole == role {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				RespondWithError(w, http.StatusForbidden, "Insufficient permissions", nil)
				return
			}

			// User has one of the required roles, proceed
			next.ServeHTTP(w, r)
		})
	}
}
