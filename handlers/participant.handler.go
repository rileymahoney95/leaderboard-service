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

// CreateParticipantRequest represents the request payload for creating a participant
type CreateParticipantRequest struct {
	ExternalID string                 `json:"external_id,omitempty" example:"external-123"`
	Name       string                 `json:"name" validate:"required" example:"John Doe"`
	Type       string                 `json:"type" validate:"required,oneof=individual team group" example:"individual" enums:"individual,team,group"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateParticipantRequest represents the request payload for updating a participant
type UpdateParticipantRequest struct {
	ExternalID *string                 `json:"external_id,omitempty" example:"external-123"`
	Name       *string                 `json:"name,omitempty" validate:"omitempty" example:"Jane Doe"`
	Type       *string                 `json:"type,omitempty" validate:"omitempty,oneof=individual team group" example:"team" enums:"individual,team,group"`
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
}

// ParticipantResponse is used for Swagger documentation
type ParticipantResponse struct {
	ID         uuid.UUID              `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ExternalID string                 `json:"external_id,omitempty" example:"external-123"`
	Name       string                 `json:"name" example:"John Doe"`
	Type       string                 `json:"type" example:"individual"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt  time.Time              `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateParticipant creates a new participant
// @Summary Create a new participant
// @Description Create a new participant with the provided details
// @Tags participants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param participant body CreateParticipantRequest true "Participant data"
// @Success 201 {object} ParticipantResponse "Created participant"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /participants [post]
func CreateParticipant(w http.ResponseWriter, r *http.Request) {
	var req CreateParticipantRequest

	err := json.NewDecoder(r.Body).Decode(&req)
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

	participant := models.Participant{
		ExternalID: req.ExternalID,
		Name:       req.Name,
		Type:       req.Type,
		Metadata:   req.Metadata,
	}

	err = db.DB.Create(&participant).Error
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create participant", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, participant)
}

// GetParticipant retrieves a participant by ID
// @Summary Get a participant by ID
// @Description Retrieve a participant by its unique ID
// @Tags participants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Participant ID"
// @Success 200 {object} ParticipantResponse "Participant details"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Router /participants/{id} [get]
func GetParticipant(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	participantID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid participant ID", err)
		return
	}

	participant := models.Participant{}
	if err := db.DB.First(&participant, "id = ?", participantID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Participant not found", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, participant)
}

// ListParticipants returns all participants
// @Summary List all participants
// @Description Get a list of all participants
// @Tags participants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ParticipantResponse "List of participants"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Router /participants [get]
func ListParticipants(w http.ResponseWriter, r *http.Request) {
	participants := []models.Participant{}
	db.DB.Find(&participants)

	middleware.RespondWithJSON(w, http.StatusOK, participants)
}

// UpdateParticipant updates an existing participant
// @Summary Update a participant
// @Description Update an existing participant with the provided details
// @Tags participants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Participant ID"
// @Param participant body UpdateParticipantRequest true "Updated participant data"
// @Success 200 {object} ParticipantResponse "Updated participant"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /participants/{id} [put]
func UpdateParticipant(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	participantID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid participant ID", err)
		return
	}

	// Fetch existing participant
	var participant models.Participant
	if err := db.DB.First(&participant, "id = ?", participantID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Participant not found", err)
		return
	}

	var req UpdateParticipantRequest
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

	// Apply the updates to the participant
	if req.ExternalID != nil {
		participant.ExternalID = *req.ExternalID
	}
	if req.Name != nil {
		participant.Name = *req.Name
	}
	if req.Type != nil {
		participant.Type = *req.Type
	}
	if req.Metadata != nil {
		participant.Metadata = *req.Metadata
	}

	// Save the updated record
	if err := db.DB.Save(&participant).Error; err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update participant", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, participant)
}

// DeleteParticipant deletes a participant by ID
// @Summary Delete a participant
// @Description Delete a participant by its ID
// @Tags participants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Participant ID"
// @Success 204 "No content"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /participants/{id} [delete]
func DeleteParticipant(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	participantID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid participant ID", err)
		return
	}

	// Check if the participant exists
	participant := models.Participant{}
	if err := db.DB.First(&participant, "id = ?", participantID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Participant not found", err)
		return
	}

	// Delete the participant
	if err := db.DB.Delete(&models.Participant{}, "id = ?", participantID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete participant", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
