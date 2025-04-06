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

// CreateLeaderboardMetricRequest represents the request payload for creating a leaderboard metric
type CreateLeaderboardMetricRequest struct {
	LeaderboardID   string  `json:"leaderboard_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	MetricID        string  `json:"metric_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440003"`
	Weight          float64 `json:"weight" validate:"required,min=0" example:"1.0"`
	DisplayPriority int     `json:"display_priority" validate:"omitempty,min=0" example:"0"`
}

// UpdateLeaderboardMetricRequest represents the request payload for updating a leaderboard metric
type UpdateLeaderboardMetricRequest struct {
	Weight          *float64 `json:"weight,omitempty" validate:"omitempty,min=0" example:"2.5"`
	DisplayPriority *int     `json:"display_priority,omitempty" validate:"omitempty,min=0" example:"1"`
}

// LeaderboardMetricResponse is used for Swagger documentation
type LeaderboardMetricResponse struct {
	ID              uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440004"`
	LeaderboardID   uuid.UUID `json:"leaderboard_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	MetricID        uuid.UUID `json:"metric_id" example:"550e8400-e29b-41d4-a716-446655440003"`
	Weight          float64   `json:"weight" example:"1.0"`
	DisplayPriority int       `json:"display_priority" example:"0"`
	CreatedAt       time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt       time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateLeaderboardMetric creates a new leaderboard metric
// @Summary Create a new leaderboard metric
// @Description Create a new metric for a leaderboard
// @Tags leaderboard-metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param leaderboard_id path string false "Leaderboard ID"
// @Param metric body CreateLeaderboardMetricRequest true "Leaderboard metric data"
// @Success 201 {object} LeaderboardMetricResponse "Created leaderboard metric"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Leaderboard or metric not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /leaderboard-metrics [post]
// @Router /leaderboards/{leaderboard_id}/metrics [post]
func CreateLeaderboardMetric(w http.ResponseWriter, r *http.Request) {
	var req CreateLeaderboardMetricRequest

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

	metricID, err := uuid.Parse(req.MetricID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid metric ID format", err)
		return
	}

	// Verify leaderboard exists
	var leaderboard models.Leaderboard
	if err := db.DB.First(&leaderboard, "id = ?", leaderboardID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard not found", err)
		return
	}

	// Set default value for display priority if not provided
	displayPriority := req.DisplayPriority
	if displayPriority < 0 {
		displayPriority = 0
	}

	leaderboardMetric := models.LeaderboardMetric{
		LeaderboardID:   leaderboardID,
		MetricID:        metricID,
		Weight:          req.Weight,
		DisplayPriority: displayPriority,
	}

	err = db.DB.Create(&leaderboardMetric).Error
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create leaderboard metric", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, leaderboardMetric)
}

// GetLeaderboardMetric retrieves a leaderboard metric by ID
// @Summary Get a leaderboard metric by ID
// @Description Retrieve a leaderboard metric by its unique ID
// @Tags leaderboard-metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Leaderboard Metric ID"
// @Success 200 {object} LeaderboardMetricResponse "Leaderboard metric details"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Router /leaderboard-metrics/{id} [get]
func GetLeaderboardMetric(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	metricID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard metric ID", err)
		return
	}

	metric := models.LeaderboardMetric{}
	if err := db.DB.First(&metric, "id = ?", metricID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard metric not found", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, metric)
}

// ListLeaderboardMetrics returns all metrics for a specific leaderboard
// @Summary List all metrics for a leaderboard
// @Description Get a list of all metrics associated with a specific leaderboard
// @Tags leaderboard-metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param leaderboard_id path string false "Filter by leaderboard ID"
// @Success 200 {array} LeaderboardMetricResponse "List of leaderboard metrics"
// @Failure 400 {object} middleware.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Router /leaderboard-metrics [get]
// @Router /leaderboards/{leaderboard_id}/metrics [get]
func ListLeaderboardMetrics(w http.ResponseWriter, r *http.Request) {
	// Check if this is a nested route call
	leaderboardIDParam := chi.URLParam(r, "id")

	// If not from nested route, check query parameter
	if leaderboardIDParam == "" {
		leaderboardIDParam = r.URL.Query().Get("leaderboard_id")
	}

	metrics := []models.LeaderboardMetric{}
	query := db.DB

	// Apply filter if provided
	if leaderboardIDParam != "" {
		leaderboardID, err := uuid.Parse(leaderboardIDParam)
		if err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard ID format", err)
			return
		}
		query = query.Where("leaderboard_id = ?", leaderboardID)
	}

	// Order by display priority
	query.Order("display_priority asc").Find(&metrics)

	middleware.RespondWithJSON(w, http.StatusOK, metrics)
}

// UpdateLeaderboardMetric updates an existing leaderboard metric
// @Summary Update a leaderboard metric
// @Description Update an existing leaderboard metric with the provided details
// @Tags leaderboard-metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Leaderboard Metric ID"
// @Param metric body UpdateLeaderboardMetricRequest true "Updated leaderboard metric data"
// @Success 200 {object} LeaderboardMetricResponse "Updated leaderboard metric"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /leaderboard-metrics/{id} [put]
func UpdateLeaderboardMetric(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	metricID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard metric ID", err)
		return
	}

	// Fetch existing metric
	var metric models.LeaderboardMetric
	if err := db.DB.First(&metric, "id = ?", metricID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard metric not found", err)
		return
	}

	var req UpdateLeaderboardMetricRequest
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

	// Apply the updates to the metric
	if req.Weight != nil {
		metric.Weight = *req.Weight
	}
	if req.DisplayPriority != nil {
		metric.DisplayPriority = *req.DisplayPriority
	}

	// Save the updated record
	if err := db.DB.Save(&metric).Error; err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update leaderboard metric", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, metric)
}

// DeleteLeaderboardMetric deletes a leaderboard metric by ID
// @Summary Delete a leaderboard metric
// @Description Delete a leaderboard metric by its ID
// @Tags leaderboard-metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Leaderboard Metric ID"
// @Success 204 "No content"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /leaderboard-metrics/{id} [delete]
func DeleteLeaderboardMetric(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	metricID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid leaderboard metric ID", err)
		return
	}

	// Check if the metric exists
	metric := models.LeaderboardMetric{}
	if err := db.DB.First(&metric, "id = ?", metricID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Leaderboard metric not found", err)
		return
	}

	// Delete the metric
	if err := db.DB.Delete(&models.LeaderboardMetric{}, "id = ?", metricID).Error; err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete leaderboard metric", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
