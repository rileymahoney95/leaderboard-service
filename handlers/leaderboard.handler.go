package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"leaderboard-service/db"
	"leaderboard-service/enums"
	"leaderboard-service/middleware"
	"leaderboard-service/models"
	"leaderboard-service/utils"
	"leaderboard-service/validation"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateLeaderboardRequest represents the request payload for creating a leaderboard
type CreateLeaderboardRequest struct {
	Name            string  `json:"name" validate:"required" example:"Weekly Tournament"`
	Description     string  `json:"description" example:"Weekly tournament for active players"`
	Category        string  `json:"category" validate:"required" example:"tournament"`
	Type            string  `json:"type" validate:"required,oneof=individual team" example:"individual" enums:"individual,team"`
	TimeFrame       string  `json:"time_frame" validate:"required,oneof=daily weekly monthly yearly all-time custom,custom_timeframe" example:"weekly" enums:"daily,weekly,monthly,yearly,all-time,custom"`
	StartDate       *string `json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z" example:"2023-01-01T00:00:00Z"`
	EndDate         *string `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z" example:"2023-01-07T23:59:59Z"`
	SortOrder       string  `json:"sort_order" validate:"required,oneof=ascending descending" example:"descending" enums:"ascending,descending"`
	VisibilityScope string  `json:"visibility_scope" validate:"required,oneof=public private" example:"public" enums:"public,private"`
	IsActive        bool    `json:"is_active" example:"true"`
	MaxEntries      int     `json:"max_entries" validate:"omitempty,min=1" example:"100"`
}

// UpdateLeaderboardRequest represents the request payload for updating a leaderboard
type UpdateLeaderboardRequest struct {
	Name            *string `json:"name,omitempty" validate:"omitempty" example:"Updated Tournament"`
	Description     *string `json:"description,omitempty" example:"Updated description"`
	Category        *string `json:"category,omitempty" validate:"omitempty" example:"competition"`
	Type            *string `json:"type,omitempty" validate:"omitempty,oneof=individual team" example:"team" enums:"individual,team"`
	TimeFrame       *string `json:"time_frame,omitempty" validate:"omitempty,oneof=daily weekly monthly yearly all-time custom,custom_timeframe" example:"monthly" enums:"daily,weekly,monthly,yearly,all-time,custom"`
	StartDate       *string `json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z" example:"2023-02-01T00:00:00Z"`
	EndDate         *string `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z" example:"2023-02-28T23:59:59Z"`
	SortOrder       *string `json:"sort_order,omitempty" validate:"omitempty,oneof=ascending descending" example:"ascending" enums:"ascending,descending"`
	VisibilityScope *string `json:"visibility_scope,omitempty" validate:"omitempty,oneof=public private" example:"private" enums:"public,private"`
	IsActive        *bool   `json:"is_active,omitempty" example:"false"`
	MaxEntries      *int    `json:"max_entries,omitempty" validate:"omitempty,min=1" example:"50"`
}

// LeaderboardResponse is used for Swagger documentation
type LeaderboardResponse struct {
	ID              uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name            string    `json:"name" example:"Weekly Tournament"`
	Description     string    `json:"description" example:"Weekly tournament for active players"`
	Category        string    `json:"category" example:"tournament"`
	Type            string    `json:"type" example:"individual"`
	TimeFrame       string    `json:"time_frame" example:"weekly"`
	StartDate       time.Time `json:"start_date,omitempty" example:"2023-01-01T00:00:00Z"`
	EndDate         time.Time `json:"end_date,omitempty" example:"2023-01-07T23:59:59Z"`
	SortOrder       string    `json:"sort_order" example:"descending"`
	VisibilityScope string    `json:"visibility_scope" example:"public"`
	IsActive        bool      `json:"is_active" example:"true"`
	MaxEntries      int       `json:"max_entries" example:"100"`
	CreatedAt       time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt       time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateLeaderboard creates a new leaderboard
// @Summary Create a new leaderboard
// @Description Create a new leaderboard with the provided details
// @Tags leaderboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param leaderboard body CreateLeaderboardRequest true "Leaderboard data"
// @Success 201 {object} LeaderboardResponse "Created leaderboard"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /leaderboards [post]
func CreateLeaderboard(w http.ResponseWriter, r *http.Request) {
	var req CreateLeaderboardRequest

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

	// Parse optional dates
	startDate, endDate := utils.ValidateDates(req.StartDate, req.EndDate)

	leaderboard := models.Leaderboard{
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		Type:            enums.LeaderboardType(req.Type),
		TimeFrame:       enums.TimeFrame(req.TimeFrame),
		StartDate:       startDate,
		EndDate:         endDate,
		SortOrder:       enums.SortOrder(req.SortOrder),
		VisibilityScope: enums.VisibilityScope(req.VisibilityScope),
		MaxEntries:      req.MaxEntries,
		IsActive:        req.IsActive,
	}

	err = db.DB.Create(&leaderboard).Error
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create leaderboard", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, leaderboard)
}

// GetLeaderboard retrieves a leaderboard by ID
// @Summary Get a leaderboard by ID
// @Description Retrieve a leaderboard by its unique ID
// @Tags leaderboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Leaderboard ID"
// @Success 200 {object} LeaderboardResponse "Leaderboard details"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Router /leaderboards/{id} [get]
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	leaderboardId, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard ID", err)
		return
	}

	leaderboard := models.Leaderboard{}
	db.DB.First(&leaderboard, "id = ?", leaderboardId)

	if leaderboard.ID == uuid.Nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard not found", nil)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, leaderboard)
}

// ListLeaderboards returns all leaderboards
// @Summary List all leaderboards
// @Description Get a list of all leaderboards
// @Tags leaderboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} LeaderboardResponse "List of leaderboards"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Router /leaderboards [get]
func ListLeaderboards(w http.ResponseWriter, r *http.Request) {
	leaderboards := []models.Leaderboard{}
	db.DB.Find(&leaderboards)

	middleware.RespondWithJSON(w, http.StatusOK, leaderboards)
}

// UpdateLeaderboard updates an existing leaderboard
// @Summary Update a leaderboard
// @Description Update an existing leaderboard with the provided details
// @Tags leaderboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Leaderboard ID"
// @Param leaderboard body UpdateLeaderboardRequest true "Updated leaderboard data"
// @Success 200 {object} LeaderboardResponse "Updated leaderboard"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /leaderboards/{id} [put]
func UpdateLeaderboard(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	leaderboardID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard ID", err)
		return
	}

	// Fetch existing leaderboard
	var leaderboard models.Leaderboard
	if err := db.DB.First(&leaderboard, "id = ?", leaderboardID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard not found", err)
		return
	}

	var req UpdateLeaderboardRequest
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

	// Apply the updates to the leaderboard
	if req.Name != nil {
		leaderboard.Name = *req.Name
	}
	if req.Description != nil {
		leaderboard.Description = *req.Description
	}
	if req.Category != nil {
		leaderboard.Category = *req.Category
	}
	if req.Type != nil {
		leaderboard.Type = enums.LeaderboardType(*req.Type)
	}
	if req.TimeFrame != nil {
		leaderboard.TimeFrame = enums.TimeFrame(*req.TimeFrame)
	}
	if req.StartDate != nil || req.EndDate != nil {
		startDate, endDate := utils.ValidateDates(req.StartDate, req.EndDate)
		if req.StartDate != nil {
			leaderboard.StartDate = startDate
		}
		if req.EndDate != nil {
			leaderboard.EndDate = endDate
		}
	}
	if req.SortOrder != nil {
		leaderboard.SortOrder = enums.SortOrder(*req.SortOrder)
	}
	if req.VisibilityScope != nil {
		leaderboard.VisibilityScope = enums.VisibilityScope(*req.VisibilityScope)
	}
	if req.MaxEntries != nil {
		leaderboard.MaxEntries = *req.MaxEntries
	}
	if req.IsActive != nil {
		leaderboard.IsActive = *req.IsActive
	}

	// Save the updated record
	if err := db.DB.Save(&leaderboard).Error; err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update leaderboard", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, leaderboard)
}

// DeleteLeaderboard deletes a leaderboard by ID
// @Summary Delete a leaderboard
// @Description Delete a leaderboard by its ID
// @Tags leaderboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Leaderboard ID"
// @Success 204 "No content"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /leaderboards/{id} [delete]
func DeleteLeaderboard(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	leaderboardID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard ID", err)
		return
	}

	leaderboard := models.Leaderboard{}
	if err := db.DB.First(&leaderboard, "id = ?", leaderboardID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard not found", err)
		return
	}

	if err := db.DB.Delete(&models.Leaderboard{}, "id = ?", leaderboardID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete leaderboard", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
