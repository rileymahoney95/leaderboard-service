package handlers

import (
	"encoding/json"
	"net/http"

	"leaderboard-service/middleware"

	"github.com/google/uuid"
)

// LoginRequest represents the login credentials
type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"password123"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType string `json:"token_type" example:"Bearer"`
	UserID    string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role      string `json:"role" example:"admin"`
}

// RegisterRequest represents registration input
type RegisterRequest struct {
	Username string `json:"username" example:"newuser"`
	Password string `json:"password" example:"securepass123"`
	Email    string `json:"email" example:"user@example.com"`
}

// Login handles user authentication and token generation
// @Summary Log in a user
// @Description Authenticate a user and generate a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// TODO: Replace with actual user authentication logic
	// For demonstration purposes, we'll accept any username/password and
	// generate a token with a random UUID as the user ID

	// In a real application, you would:
	// 1. Validate username/password against the database
	// 2. If valid, generate a token with the user's actual ID and role
	// 3. Return the token to the client

	// Mock user authentication
	if req.Username == "" || req.Password == "" {
		middleware.RespondWithError(w, http.StatusBadRequest, "Username and password are required", nil)
		return
	}

	// Generate a mock user ID and role
	userID := uuid.New().String()
	userRole := "user"

	// Assign admin role for a specific username (for testing)
	if req.Username == "admin" {
		userRole = "admin"
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(userID, userRole)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	// Create response
	resp := LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		UserID:    userID,
		Role:      userRole,
	}

	middleware.RespondWithJSON(w, http.StatusOK, resp)
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param registerRequest body RegisterRequest true "Registration data"
// @Success 201 {object} LoginResponse "Registration successful"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement user registration
	// This is a placeholder for future implementation
	middleware.RespondWithError(w, http.StatusNotImplemented, "Registration not implemented yet", nil)
}
