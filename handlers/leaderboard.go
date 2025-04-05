package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"leaderboard-service/db"
	"leaderboard-service/enums"
	"leaderboard-service/middleware"
	"leaderboard-service/models"
	"leaderboard-service/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateLeaderboardRequest represents the request payload for creating a leaderboard
type CreateLeaderboardRequest struct {
	Name            string  `json:"name" example:"Weekly Tournament"`
	Description     string  `json:"description" example:"Weekly tournament for active players"`
	Category        string  `json:"category" example:"tournament"`
	Type            string  `json:"type" example:"score" enums:"score,time,points"`
	TimeFrame       string  `json:"time_frame" example:"weekly" enums:"daily,weekly,monthly,yearly,all_time,custom"`
	StartDate       *string `json:"start_date,omitempty" example:"2023-01-01T00:00:00Z"`
	EndDate         *string `json:"end_date,omitempty" example:"2023-01-07T23:59:59Z"`
	SortOrder       string  `json:"sort_order" example:"desc" enums:"asc,desc"`
	VisibilityScope string  `json:"visibility_scope" example:"public" enums:"public,private,restricted"`
	IsActive        bool    `json:"is_active" example:"true"`
	MaxEntries      int     `json:"max_entries" example:"100"`
}

// UpdateLeaderboardRequest represents the request payload for updating a leaderboard
type UpdateLeaderboardRequest struct {
	Name            *string `json:"name,omitempty" example:"Updated Tournament"`
	Description     *string `json:"description,omitempty" example:"Updated description"`
	Category        *string `json:"category,omitempty" example:"competition"`
	Type            *string `json:"type,omitempty" example:"points" enums:"score,time,points"`
	TimeFrame       *string `json:"time_frame,omitempty" example:"monthly" enums:"daily,weekly,monthly,yearly,all_time,custom"`
	StartDate       *string `json:"start_date,omitempty" example:"2023-02-01T00:00:00Z"`
	EndDate         *string `json:"end_date,omitempty" example:"2023-02-28T23:59:59Z"`
	SortOrder       *string `json:"sort_order,omitempty" example:"asc" enums:"asc,desc"`
	VisibilityScope *string `json:"visibility_scope,omitempty" example:"restricted" enums:"public,private,restricted"`
	IsActive        *bool   `json:"is_active,omitempty" example:"false"`
	MaxEntries      *int    `json:"max_entries,omitempty" example:"50"`
}

// LeaderboardResponse is used for Swagger documentation
type LeaderboardResponse struct {
	ID              uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name            string    `json:"name" example:"Weekly Tournament"`
	Description     string    `json:"description" example:"Weekly tournament for active players"`
	Category        string    `json:"category" example:"tournament"`
	Type            string    `json:"type" example:"score"`
	TimeFrame       string    `json:"time_frame" example:"weekly"`
	StartDate       time.Time `json:"start_date,omitempty" example:"2023-01-01T00:00:00Z"`
	EndDate         time.Time `json:"end_date,omitempty" example:"2023-01-07T23:59:59Z"`
	SortOrder       string    `json:"sort_order" example:"desc"`
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

	// Validate required fields
	if err := validateCreateLeaderboardRequest(req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Validation error", err)
		return
	}

	// Convert string to enum and validate
	leaderboardType := enums.LeaderboardType(req.Type)
	if !leaderboardType.Valid() {
		errorMessage := "Invalid leaderboard type, valid values are: " + strings.Join(enums.GetValidLeaderboardTypes(), ", ")
		middleware.RespondWithError(w, http.StatusBadRequest, errorMessage, nil)
		return
	}

	timeFrame := enums.TimeFrame(req.TimeFrame)
	if !timeFrame.Valid() {
		errorMessage := "Invalid leaderboard time frame, valid values are: " + strings.Join(enums.GetValidTimeFrames(), ", ")
		middleware.RespondWithError(w, http.StatusBadRequest, errorMessage, nil)
		return
	}

	sortOrder := enums.SortOrder(req.SortOrder)
	if !sortOrder.Valid() {
		errorMessage := "Invalid leaderboard sort order, valid values are: " + strings.Join(enums.GetValidSortOrders(), ", ")
		middleware.RespondWithError(w, http.StatusBadRequest, errorMessage, nil)
		return
	}

	visibilityScope := enums.VisibilityScope(req.VisibilityScope)
	if !visibilityScope.Valid() {
		errorMessage := "Invalid leaderboard visibility scope, valid values are: " + strings.Join(enums.GetValidVisibilityScopes(), ", ")
		middleware.RespondWithError(w, http.StatusBadRequest, errorMessage, nil)
		return
	}

	// Parse optional dates
	startDate, endDate := utils.ValidateDates(req.StartDate, req.EndDate)

	leaderboard := models.Leaderboard{
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		Type:            leaderboardType,
		TimeFrame:       timeFrame,
		StartDate:       startDate,
		EndDate:         endDate,
		SortOrder:       sortOrder,
		VisibilityScope: visibilityScope,
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

	// Read the raw request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Failed to read request body", err)
		return
	}

	// Fetch existing leaderboard
	var leaderboard models.Leaderboard
	if err := db.DB.First(&leaderboard, "id = ?", leaderboardID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard not found", err)
		return
	}

	// Parse the JSON request
	var rawData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &rawData); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid JSON in request body", err)
		return
	}

	// Handle specific fields with case-insensitive matching
	var req UpdateLeaderboardRequest
	processRequestFields(rawData, &req)

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
		leaderboardType := enums.LeaderboardType(*req.Type)
		if !leaderboardType.Valid() {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard type", nil)
			return
		}
		leaderboard.Type = leaderboardType
	}
	if req.TimeFrame != nil {
		timeFrame := enums.TimeFrame(*req.TimeFrame)
		if !timeFrame.Valid() {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard time frame", nil)
			return
		}
		leaderboard.TimeFrame = timeFrame
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
		sortOrder := enums.SortOrder(*req.SortOrder)
		if !sortOrder.Valid() {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard sort order", nil)
			return
		}
		leaderboard.SortOrder = sortOrder
	}
	if req.VisibilityScope != nil {
		visibilityScope := enums.VisibilityScope(*req.VisibilityScope)
		if !visibilityScope.Valid() {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard visibility scope", nil)
			return
		}
		leaderboard.VisibilityScope = visibilityScope
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

/*
* Validate the request body for creating a leaderboard
 */
func validateCreateLeaderboardRequest(req CreateLeaderboardRequest) error {
	missingFields := []string{}

	// Check required fields based on the model's not null constraints
	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}
	if req.Category == "" {
		missingFields = append(missingFields, "category")
	}
	if req.Type == "" {
		missingFields = append(missingFields, "type")
	}
	if req.TimeFrame == "" {
		missingFields = append(missingFields, "time_frame")
	}
	if req.SortOrder == "" {
		missingFields = append(missingFields, "sort_order")
	}
	if req.VisibilityScope == "" {
		missingFields = append(missingFields, "visibility_scope")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

/*
* Process JSON fields with case-insensitive matching
 */
func processRequestFields(rawData map[string]interface{}, req *UpdateLeaderboardRequest) {
	for key, value := range rawData {
		lowerKey := strings.ToLower(key)

		switch lowerKey {
		case "name":
			if str, ok := value.(string); ok {
				req.Name = &str
			}
		case "description":
			if str, ok := value.(string); ok {
				req.Description = &str
			}
		case "category":
			if str, ok := value.(string); ok {
				req.Category = &str
			}
		case "type":
			if str, ok := value.(string); ok {
				req.Type = &str
			}
		case "timeframe", "time_frame":
			if str, ok := value.(string); ok {
				req.TimeFrame = &str
			}
		case "startdate", "start_date":
			if str, ok := value.(string); ok {
				req.StartDate = &str
			}
		case "enddate", "end_date":
			if str, ok := value.(string); ok {
				req.EndDate = &str
			}
		case "sortorder", "sort_order":
			if str, ok := value.(string); ok {
				req.SortOrder = &str
			}
		case "visibilityscope", "visibility_scope":
			if str, ok := value.(string); ok {
				req.VisibilityScope = &str
			}
		case "isactive", "is_active":
			if b, ok := value.(bool); ok {
				req.IsActive = &b
			}
		case "maxentries", "max_entries":
			if num, ok := value.(float64); ok {
				intVal := int(num)
				req.MaxEntries = &intVal
			}
		}
	}
}
