package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestLogger logs information about each request including user information if available
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Process the request
		next.ServeHTTP(ww, r)

		// Calculate request duration
		duration := time.Since(start)

		// Get request ID if available
		requestID := middleware.GetReqID(r.Context())

		// Try to get user information if available
		userInfo := "anonymous"
		if claims, err := GetUserFromContext(r.Context()); err == nil {
			userInfo = fmt.Sprintf("user_id=%s role=%s", claims.UserID, claims.Role)
		}

		// Log the request details
		log.Printf(
			"[%s] %s %s %s - %d %s - %s",
			requestID,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			ww.Status(),
			duration,
			userInfo,
		)
	})
}
