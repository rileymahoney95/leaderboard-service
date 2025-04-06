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

// CreateMetricValueRequest represents the request payload for creating a metric value
type CreateMetricValueRequest struct {
	MetricID      string      `json:"metric_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ParticipantID string      `json:"participant_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440001"`
	Value         float64     `json:"value" validate:"required" example:"42.5"`
	Timestamp     *time.Time  `json:"timestamp,omitempty" example:"2023-01-01T00:00:00Z"`
	Source        string      `json:"source,omitempty" example:"call_system"`
	Context       interface{} `json:"context,omitempty"`
}

// UpdateMetricValueRequest represents the request payload for updating a metric value
type UpdateMetricValueRequest struct {
	Value     *float64     `json:"value,omitempty" validate:"omitempty" example:"50.75"`
	Timestamp *time.Time   `json:"timestamp,omitempty" example:"2023-01-02T00:00:00Z"`
	Source    *string      `json:"source,omitempty" example:"text_system"`
	Context   *interface{} `json:"context,omitempty"`
}

// MetricValueResponse is used for Swagger documentation
type MetricValueResponse struct {
	ID            uuid.UUID   `json:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	MetricID      uuid.UUID   `json:"metric_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ParticipantID uuid.UUID   `json:"participant_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Value         float64     `json:"value" example:"42.5"`
	Timestamp     time.Time   `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Source        string      `json:"source,omitempty" example:"call_system"`
	Context       interface{} `json:"context,omitempty"`
	CreatedAt     time.Time   `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time   `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

type MetricValueHandler struct {
	service services.MetricValueService
}

func NewMetricValueHandler() *MetricValueHandler {
	metricValueRepo := repositories.NewMetricValueRepository()
	metricRepo := repositories.NewMetricRepository()
	participantRepo := repositories.NewParticipantRepository()
	service := services.NewMetricValueService(metricValueRepo, metricRepo, participantRepo)

	return &MetricValueHandler{
		service: service,
	}
}

// CreateMetricValue creates a new metric value
// @Summary Create a new metric value
// @Description Create a new metric value record for a participant
// @Tags metric-values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param metric_id path string false "Metric ID"
// @Param participant_id path string false "Participant ID"
// @Param metric_value body CreateMetricValueRequest true "Metric value data"
// @Success 201 {object} MetricValueResponse "Created metric value"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Metric or participant not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /metric-values [post]
// @Router /metrics/{metric_id}/values [post]
// @Router /participants/{participant_id}/metric-values [post]
func (h *MetricValueHandler) CreateMetricValue(w http.ResponseWriter, r *http.Request) {
	var req CreateMetricValueRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Get path parameters from nested routes
	metricIDPath := chi.URLParam(r, "id")
	participantIDPath := chi.URLParam(r, "id")

	// Determine the context of the call (which nested route we're using)
	routePath := r.URL.Path
	isMetricNested := false
	isParticipantNested := false

	if len(routePath) >= 8 && routePath[:8] == "/metrics" {
		isMetricNested = true
	} else if len(routePath) >= 13 && routePath[:13] == "/participants" {
		isParticipantNested = true
	}

	// Override request values with path parameters if available
	if isMetricNested && metricIDPath != "" {
		req.MetricID = metricIDPath
	}

	if isParticipantNested && participantIDPath != "" {
		req.ParticipantID = participantIDPath
	}

	// Validate using validator package
	if err := validation.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		middleware.RespondWithError(w, http.StatusBadRequest, "Validation error", validation.FormatValidationErrors(validationErrors))
		return
	}

	// Parse UUIDs
	metricID, err := uuid.Parse(req.MetricID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid metric ID format", err)
		return
	}

	participantID, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid participant ID format", err)
		return
	}

	// Set timestamp to current time if not provided
	timestamp := time.Now()
	if req.Timestamp != nil {
		timestamp = *req.Timestamp
	}

	metricValue, err := h.service.CreateMetricValue(
		metricID,
		participantID,
		req.Value,
		timestamp,
		req.Source,
		req.Context,
	)

	if err != nil {
		if err.Error() == "metric not found" || err.Error() == "participant not found" {
			middleware.RespondWithError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create metric value", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, metricValue)
}

// GetMetricValue retrieves a metric value by ID
// @Summary Get a metric value by ID
// @Description Retrieve a metric value by its unique ID
// @Tags metric-values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Metric Value ID"
// @Success 200 {object} MetricValueResponse "Metric value details"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Router /metric-values/{id} [get]
func (h *MetricValueHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	valueID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid metric value ID", err)
		return
	}

	value, err := h.service.GetMetricValue(valueID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Metric value not found", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, value)
}

// ListMetricValues returns metric values with optional filtering
// @Summary List metric values
// @Description Get a list of metric values with optional filtering by metric ID and/or participant ID
// @Tags metric-values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param metric_id path string false "Filter by metric ID"
// @Param participant_id path string false "Filter by participant ID"
// @Param from_time query string false "Filter by timestamp (greater than or equal)" format(date-time)
// @Param to_time query string false "Filter by timestamp (less than or equal)" format(date-time)
// @Success 200 {array} MetricValueResponse "List of metric values"
// @Failure 400 {object} middleware.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Router /metric-values [get]
// @Router /metrics/{metric_id}/values [get]
// @Router /participants/{participant_id}/metric-values [get]
func (h *MetricValueHandler) ListMetricValues(w http.ResponseWriter, r *http.Request) {
	// Get path parameters from nested routes
	metricIDPath := chi.URLParam(r, "id")
	participantIDPath := chi.URLParam(r, "id")

	// Determine the context of the call (which nested route we're using)
	routePath := r.URL.Path
	isMetricNested := false
	isParticipantNested := false

	if len(routePath) >= 8 && routePath[:8] == "/metrics" {
		isMetricNested = true
	} else if len(routePath) >= 13 && routePath[:13] == "/participants" {
		isParticipantNested = true
	}

	// Get query parameters (for flat route)
	metricIDQuery := r.URL.Query().Get("metric_id")
	participantIDQuery := r.URL.Query().Get("participant_id")
	fromTimeParam := r.URL.Query().Get("from_time")
	toTimeParam := r.URL.Query().Get("to_time")

	metricIDParam := ""
	participantIDParam := ""

	// Determine which param to use based on route context
	if isMetricNested {
		metricIDParam = metricIDPath
		participantIDParam = participantIDQuery
	} else if isParticipantNested {
		metricIDParam = metricIDQuery
		participantIDParam = participantIDPath
	} else {
		// Flat route
		metricIDParam = metricIDQuery
		participantIDParam = participantIDQuery
	}

	var metricID *uuid.UUID
	var participantID *uuid.UUID
	var fromTime *time.Time
	var toTime *time.Time

	// Parse metricID if provided
	if metricIDParam != "" {
		parsedMetricID, err := uuid.Parse(metricIDParam)
		if err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid metric ID format", err)
			return
		}
		metricID = &parsedMetricID
	}

	// Parse participantID if provided
	if participantIDParam != "" {
		parsedParticipantID, err := uuid.Parse(participantIDParam)
		if err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid participant ID format", err)
			return
		}
		participantID = &parsedParticipantID
	}

	// Parse fromTime if provided
	if fromTimeParam != "" {
		parsedFromTime, err := time.Parse(time.RFC3339, fromTimeParam)
		if err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid from_time format, use RFC3339", err)
			return
		}
		fromTime = &parsedFromTime
	}

	// Parse toTime if provided
	if toTimeParam != "" {
		parsedToTime, err := time.Parse(time.RFC3339, toTimeParam)
		if err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid to_time format, use RFC3339", err)
			return
		}
		toTime = &parsedToTime
	}

	values, err := h.service.ListFilteredMetricValues(metricID, participantID, fromTime, toTime)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch metric values", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, values)
}

// UpdateMetricValue updates an existing metric value
// @Summary Update a metric value
// @Description Update an existing metric value with the provided details
// @Tags metric-values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Metric Value ID"
// @Param metric_value body UpdateMetricValueRequest true "Updated metric value data"
// @Success 200 {object} MetricValueResponse "Updated metric value"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /metric-values/{id} [put]
func (h *MetricValueHandler) UpdateMetricValue(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	valueID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid metric value ID", err)
		return
	}

	var req UpdateMetricValueRequest
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

	updatedValue, err := h.service.UpdateMetricValue(
		valueID,
		req.Value,
		req.Timestamp,
		req.Source,
		req.Context,
	)

	if err != nil {
		if err.Error() == "metric value not found" {
			middleware.RespondWithError(w, http.StatusNotFound, "Metric value not found", err)
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update metric value", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, updatedValue)
}

// DeleteMetricValue deletes a metric value by ID
// @Summary Delete a metric value
// @Description Delete a metric value by its ID
// @Tags metric-values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Metric Value ID"
// @Success 204 "No content"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /metric-values/{id} [delete]
func (h *MetricValueHandler) DeleteMetricValue(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	valueID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid metric value ID", err)
		return
	}

	err = h.service.DeleteMetricValue(valueID)
	if err != nil {
		if err.Error() == "metric value not found" {
			middleware.RespondWithError(w, http.StatusNotFound, "Metric value not found", err)
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete metric value", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
