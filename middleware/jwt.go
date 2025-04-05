package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Define custom claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// Define error constants
var (
	ErrTokenMissing      = errors.New("token is missing")
	ErrTokenInvalid      = errors.New("token is invalid")
	ErrTokenExpired      = errors.New("token is expired")
	ErrInvalidSignMethod = errors.New("invalid signing method")
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserContextKey is the key used to store user information in the request context
	UserContextKey ContextKey = "user"
)

// JWTAuth middleware authenticates requests using JWT
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, ErrTokenMissing.Error(), http.StatusUnauthorized)
			return
		}

		// Extract token from header
		tokenString := extractTokenFromHeader(authHeader)
		if tokenString == "" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Parse and validate the token
		claims, err := validateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Add claims to the request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractTokenFromHeader extracts the JWT token from various header formats
func extractTokenFromHeader(authHeader string) string {
	// Check standard Bearer format first
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Check for JWT directly (might be the case with Swagger UI)
	if strings.HasPrefix(authHeader, "JWT ") {
		return strings.TrimPrefix(authHeader, "JWT ")
	}

	// Check if it's just the token without any prefix
	if len(authHeader) > 20 && !strings.Contains(authHeader, " ") {
		return authHeader
	}

	// Check if token is after the keyword 'Bearer' without having a space
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer") {
		return strings.TrimPrefix(authHeader, "bearer")
	}

	return ""
}

// validateToken parses and validates the JWT token
func validateToken(tokenString string) (*Claims, error) {
	// Get JWT secret from environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET environment variable not set")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: %v", ErrInvalidSignMethod, token.Header["alg"])
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}

	// Extract the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID, role string) (string, error) {
	// Get JWT secret and expiration from environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET environment variable not set")
	}

	expirationHours := 24 // Default to 24 hours
	if os.Getenv("JWT_EXPIRATION_HOURS") != "" {
		fmt.Sscanf(os.Getenv("JWT_EXPIRATION_HOURS"), "%d", &expirationHours)
	}

	// Set up the claims
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "leaderboard-service",
			Subject:   userID,
		},
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserFromContext extracts user claims from the request context
func GetUserFromContext(ctx context.Context) (*Claims, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	// Extract user claims from context
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	if !ok || claims == nil {
		return nil, errors.New("user not found in context")
	}

	return claims, nil
}
