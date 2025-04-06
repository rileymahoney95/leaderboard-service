package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"leaderboard-service/middleware"
	"leaderboard-service/repositories"
	"leaderboard-service/services"
	"leaderboard-service/validation"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateLeaderboardEntryRequest represents the request payload for creating a leaderboard entry
type CreateLeaderboardEntryRequest struct {
	LeaderboardID string    `json:"leaderboard_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ParticipantID string    `json:"participant_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440001"`
	Score         float64   `json:"score" validate:"required" example:"100.5"`
	Rank          int       `json:"rank" validate:"required,min=1" example:"1"`
	LastUpdated   time.Time `json:"last_updated,omitempty" example:"2023-01-01T00:00:00Z"`
}

// UpdateLeaderboardEntryRequest represents the request payload for updating a leaderboard entry
type UpdateLeaderboardEntryRequest struct {
	Score       *float64   `json:"score,omitempty" validate:"omitempty" example:"200.75"`
	Rank        *int       `json:"rank,omitempty" validate:"omitempty,min=1" example:"2"`
	LastUpdated *time.Time `json:"last_updated,omitempty" example:"2023-01-02T00:00:00Z"`
}

// LeaderboardEntryResponse is used for Swagger documentation
type LeaderboardEntryResponse struct {
	ID            uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	LeaderboardID uuid.UUID `json:"leaderboard_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ParticipantID uuid.UUID `json:"participant_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Rank          int       `json:"rank" example:"1"`
	Score         float64   `json:"score" example:"100.5"`
	LastUpdated   time.Time `json:"last_updated" example:"2023-01-01T00:00:00Z"`
	CreatedAt     time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

type LeaderboardEntryHandler struct {
	service services.LeaderboardEntryService
}

func NewLeaderboardEntryHandler() *LeaderboardEntryHandler {
	leaderboardEntryRepo := repositories.NewLeaderboardEntryRepository()
	leaderboardRepo := repositories.NewLeaderboardRepository()
	participantRepo := repositories.NewParticipantRepository()
	service := services.NewLeaderboardEntryService(leaderboardEntryRepo, leaderboardRepo, participantRepo)

	return &LeaderboardEntryHandler{
		service: service,
	}
}

// CreateLeaderboardEntry creates a new leaderboard entry
// @Summary Create a new leaderboard entry
// @Description Create a new entry/ranking in a leaderboard
// @Tags leaderboard-entries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param leaderboard_id path string false "Leaderboard ID"
// @Param entry body CreateLeaderboardEntryRequest true "Leaderboard entry data"
// @Success 201 {object} LeaderboardEntryResponse "Created leaderboard entry"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Leaderboard or participant not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /leaderboard-entries [post]
// @Router /leaderboards/{leaderboard_id}/entries [post]
func (h *LeaderboardEntryHandler) CreateLeaderboardEntry(w http.ResponseWriter, r *http.Request) {
	var req CreateLeaderboardEntryRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Check if this is a nested route call
	leaderboardIDPath := chi.URLParam(r, "id")

	// Override request values with path parameters if available
	if leaderboardIDPath != "" {
		req.LeaderboardID = leaderboardIDPath
	}

	// Validate using validator package
	if err := validation.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		middleware.RespondWithError(w, http.StatusBadRequest, "Validation error", validation.FormatValidationErrors(validationErrors))
		return
	}

	// Parse UUIDs
	leaderboardID, err := uuid.Parse(req.LeaderboardID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard ID format", err)
		return
	}

	participantID, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid participant ID format", err)
		return
	}

	entry, err := h.service.CreateLeaderboardEntry(
		leaderboardID,
		participantID,
		req.Score,
		req.Rank,
		req.LastUpdated,
	)

	if err != nil {
		if err.Error() == "leaderboard not found" || err.Error() == "participant not found" {
			middleware.RespondWithError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create leaderboard entry", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, entry)
}

// GetLeaderboardEntry retrieves a leaderboard entry by ID
// @Summary Get a leaderboard entry by ID
// @Description Retrieve a leaderboard entry by its unique ID
// @Tags leaderboard-entries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Leaderboard Entry ID"
// @Success 200 {object} LeaderboardEntryResponse "Leaderboard entry details"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Router /leaderboard-entries/{id} [get]
func (h *LeaderboardEntryHandler) GetLeaderboardEntry(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard entry ID", err)
		return
	}

	entry, err := h.service.GetLeaderboardEntry(entryID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard entry not found", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, entry)
}

// ListLeaderboardEntries returns all entries for a specific leaderboard
// @Summary List all entries for a leaderboard
// @Description Get a list of all entries/rankings for a specific leaderboard
// @Tags leaderboard-entries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param leaderboard_id path string false "Filter by leaderboard ID"
// @Param participant_id query string false "Filter by participant ID"
// @Success 200 {array} LeaderboardEntryResponse "List of leaderboard entries"
// @Failure 400 {object} middleware.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Router /leaderboard-entries [get]
// @Router /leaderboards/{leaderboard_id}/entries [get]
func (h *LeaderboardEntryHandler) ListLeaderboardEntries(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	participantIDParam := r.URL.Query().Get("participant_id")

	// Check if this is a nested route call
	leaderboardIDParam := chi.URLParam(r, "id")

	// If not from nested route, check query parameter
	if leaderboardIDParam == "" {
		leaderboardIDParam = r.URL.Query().Get("leaderboard_id")
	}

	var leaderboardID *uuid.UUID
	var participantID *uuid.UUID

	// Parse leaderboardID if provided
	if leaderboardIDParam != "" {
		parsedID, err := uuid.Parse(leaderboardIDParam)
		if err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard ID format", err)
			return
		}
		leaderboardID = &parsedID
	}

	// Parse participantID if provided
	if participantIDParam != "" {
		parsedID, err := uuid.Parse(participantIDParam)
		if err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid participant ID format", err)
			return
		}
		participantID = &parsedID
	}

	entries, err := h.service.ListFilteredLeaderboardEntries(leaderboardID, participantID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch leaderboard entries", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, entries)
}

// UpdateLeaderboardEntry updates an existing leaderboard entry
// @Summary Update a leaderboard entry
// @Description Update an existing leaderboard entry with the provided details
// @Tags leaderboard-entries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Leaderboard Entry ID"
// @Param entry body UpdateLeaderboardEntryRequest true "Updated leaderboard entry data"
// @Success 200 {object} LeaderboardEntryResponse "Updated leaderboard entry"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /leaderboard-entries/{id} [put]
func (h *LeaderboardEntryHandler) UpdateLeaderboardEntry(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard entry ID", err)
		return
	}

	var req UpdateLeaderboardEntryRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Validate using validator package
	if err := validation.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		middleware.RespondWithError(w, http.StatusBadRequest, "Validation error", validation.FormatValidationErrors(validationErrors))
		return
	}

	updatedEntry, err := h.service.UpdateLeaderboardEntry(
		entryID,
		req.Score,
		req.Rank,
		req.LastUpdated,
	)

	if err != nil {
		if err.Error() == "leaderboard entry not found" {
			middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard entry not found", err)
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update leaderboard entry", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, updatedEntry)
}

// DeleteLeaderboardEntry deletes a leaderboard entry by ID
// @Summary Delete a leaderboard entry
// @Description Delete a leaderboard entry by its ID
// @Tags leaderboard-entries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Leaderboard Entry ID"
// @Success 204 "No content"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /leaderboard-entries/{id} [delete]
func (h *LeaderboardEntryHandler) DeleteLeaderboardEntry(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard entry ID", err)
		return
	}

	err = h.service.DeleteLeaderboardEntry(entryID)
	if err != nil {
		if err.Error() == "leaderboard entry not found" {
			middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard entry not found", err)
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete leaderboard entry", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
