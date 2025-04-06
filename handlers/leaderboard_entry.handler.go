package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"leaderboard-service/db"
	"leaderboard-service/middleware"
	"leaderboard-service/models"
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
func CreateLeaderboardEntry(w http.ResponseWriter, r *http.Request) {
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

	// Verify leaderboard exists
	var leaderboard models.Leaderboard
	if err := db.DB.First(&leaderboard, "id = ?", leaderboardID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard not found", err)
		return
	}

	// Verify participant exists
	var participant models.Participant
	if err := db.DB.First(&participant, "id = ?", participantID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Participant not found", err)
		return
	}

	// Set last updated to current time if not provided
	lastUpdated := req.LastUpdated
	if lastUpdated.IsZero() {
		lastUpdated = time.Now()
	}

	leaderboardEntry := models.LeaderboardEntry{
		LeaderboardID: leaderboardID,
		ParticipantID: participantID,
		Rank:          req.Rank,
		Score:         req.Score,
		LastUpdated:   lastUpdated,
	}

	err = db.DB.Create(&leaderboardEntry).Error
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create leaderboard entry", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, leaderboardEntry)
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
func GetLeaderboardEntry(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard entry ID", err)
		return
	}

	entry := models.LeaderboardEntry{}
	if err := db.DB.First(&entry, "id = ?", entryID).Error; err != nil {
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
func ListLeaderboardEntries(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	participantIDParam := r.URL.Query().Get("participant_id")

	// Check if this is a nested route call
	leaderboardIDParam := chi.URLParam(r, "id")

	// If not from nested route, check query parameter
	if leaderboardIDParam == "" {
		leaderboardIDParam = r.URL.Query().Get("leaderboard_id")
	}

	entries := []models.LeaderboardEntry{}
	query := db.DB

	// Apply filters if provided
	if leaderboardIDParam != "" {
		leaderboardID, err := uuid.Parse(leaderboardIDParam)
		if err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard ID format", err)
			return
		}
		query = query.Where("leaderboard_id = ?", leaderboardID)
	}

	if participantIDParam != "" {
		participantID, err := uuid.Parse(participantIDParam)
		if err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid participant ID format", err)
			return
		}
		query = query.Where("participant_id = ?", participantID)
	}

	// Order by rank
	query.Order("rank asc").Find(&entries)

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
func UpdateLeaderboardEntry(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard entry ID", err)
		return
	}

	// Fetch existing entry
	var entry models.LeaderboardEntry
	if err := db.DB.First(&entry, "id = ?", entryID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard entry not found", err)
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

	// Apply the updates to the entry
	if req.Score != nil {
		entry.Score = *req.Score
	}
	if req.Rank != nil {
		entry.Rank = *req.Rank
	}
	if req.LastUpdated != nil {
		entry.LastUpdated = *req.LastUpdated
	} else {
		// Update the LastUpdated field if not explicitly provided
		entry.LastUpdated = time.Now()
	}

	// Save the updated record
	if err := db.DB.Save(&entry).Error; err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update leaderboard entry", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, entry)
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
func DeleteLeaderboardEntry(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard entry ID", err)
		return
	}

	// Check if the entry exists
	entry := models.LeaderboardEntry{}
	if err := db.DB.First(&entry, "id = ?", entryID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard entry not found", err)
		return
	}

	// Delete the entry
	if err := db.DB.Delete(&models.LeaderboardEntry{}, "id = ?", entryID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete leaderboard entry", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
