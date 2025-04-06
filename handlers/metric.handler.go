package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"leaderboard-service/enums"
	"leaderboard-service/middleware"
	"leaderboard-service/repositories"
	"leaderboard-service/services"
	"leaderboard-service/validation"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateMetricRequest represents the request payload for creating a metric
type CreateMetricRequest struct {
	Name            string `json:"name" validate:"required" example:"monthly_calls_completed"`
	Description     string `json:"description" example:"Number of calls completed in a month"`
	DataType        string `json:"data_type" validate:"required,oneof=integer decimal boolean string" example:"integer" enums:"integer,decimal,boolean,string"`
	Unit            string `json:"unit" example:"calls"`
	AggregationType string `json:"aggregation_type" validate:"required,oneof=sum average count min max last" example:"sum" enums:"sum,average,count,min,max,last"`
	ResetPeriod     string `json:"reset_period" validate:"required,oneof=none daily weekly monthly yearly" example:"monthly" enums:"none,daily,weekly,monthly,yearly"`
	IsHigherBetter  bool   `json:"is_higher_better" example:"true"`
}

// UpdateMetricRequest represents the request payload for updating a metric
type UpdateMetricRequest struct {
	Name            *string `json:"name,omitempty" validate:"omitempty" example:"monthly_texts_answered"`
	Description     *string `json:"description,omitempty" example:"Number of texts answered in a month"`
	DataType        *string `json:"data_type,omitempty" validate:"omitempty,oneof=integer decimal boolean string" example:"integer" enums:"integer,decimal,boolean,string"`
	Unit            *string `json:"unit,omitempty" example:"texts"`
	AggregationType *string `json:"aggregation_type,omitempty" validate:"omitempty,oneof=sum average count min max last" example:"sum" enums:"sum,average,count,min,max,last"`
	ResetPeriod     *string `json:"reset_period,omitempty" validate:"omitempty,oneof=none daily weekly monthly yearly" example:"monthly" enums:"none,daily,weekly,monthly,yearly"`
	IsHigherBetter  *bool   `json:"is_higher_better,omitempty" example:"true"`
}

// MetricResponse is used for Swagger documentation
type MetricResponse struct {
	ID              uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name            string    `json:"name" example:"monthly_calls_completed"`
	Description     string    `json:"description" example:"Number of calls completed in a month"`
	DataType        string    `json:"data_type" example:"integer"`
	Unit            string    `json:"unit" example:"calls"`
	AggregationType string    `json:"aggregation_type" example:"sum"`
	ResetPeriod     string    `json:"reset_period" example:"monthly"`
	IsHigherBetter  bool      `json:"is_higher_better" example:"true"`
	CreatedAt       time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt       time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

type MetricHandler struct {
	service services.MetricService
}

func NewMetricHandler() *MetricHandler {
	repo := repositories.NewMetricRepository()
	service := services.NewMetricService(repo)
	return &MetricHandler{
		service: service,
	}
}

// CreateMetric creates a new metric
// @Summary Create a new metric
// @Description Create a new metric with the provided details
// @Tags metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param metric body CreateMetricRequest true "Metric data"
// @Success 201 {object} MetricResponse "Created metric"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /metrics [post]
func (h *MetricHandler) CreateMetric(w http.ResponseWriter, r *http.Request) {
	var req CreateMetricRequest

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

	metric, err := h.service.CreateMetric(
		req.Name,
		req.Description,
		enums.MetricDataType(req.DataType),
		req.Unit,
		enums.AggregationType(req.AggregationType),
		enums.ResetPeriod(req.ResetPeriod),
		req.IsHigherBetter,
	)

	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create metric", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, metric)
}

// GetMetric retrieves a metric by ID
// @Summary Get a metric by ID
// @Description Retrieve a metric by its unique ID
// @Tags metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Metric ID"
// @Success 200 {object} MetricResponse "Metric details"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Router /metrics/{id} [get]
func (h *MetricHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	metricID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid metric ID", err)
		return
	}

	metric, err := h.service.GetMetric(metricID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Metric not found", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, metric)
}

// ListMetrics returns all metrics
// @Summary List all metrics
// @Description Get a list of all metrics
// @Tags metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} MetricResponse "List of metrics"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Router /metrics [get]
func (h *MetricHandler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.service.ListMetrics()
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch metrics", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, metrics)
}

// UpdateMetric updates an existing metric
// @Summary Update a metric
// @Description Update an existing metric with the provided details
// @Tags metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Metric ID"
// @Param metric body UpdateMetricRequest true "Updated metric data"
// @Success 200 {object} MetricResponse "Updated metric"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /metrics/{id} [put]
func (h *MetricHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	metricID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid metric ID", err)
		return
	}

	var req UpdateMetricRequest
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

	// Convert string types to enum types
	var dataType *enums.MetricDataType
	if req.DataType != nil {
		dt := enums.MetricDataType(*req.DataType)
		dataType = &dt
	}

	var aggregationType *enums.AggregationType
	if req.AggregationType != nil {
		at := enums.AggregationType(*req.AggregationType)
		aggregationType = &at
	}

	var resetPeriod *enums.ResetPeriod
	if req.ResetPeriod != nil {
		rp := enums.ResetPeriod(*req.ResetPeriod)
		resetPeriod = &rp
	}

	updatedMetric, err := h.service.UpdateMetric(
		metricID,
		req.Name,
		req.Description,
		dataType,
		req.Unit,
		aggregationType,
		resetPeriod,
		req.IsHigherBetter,
	)

	if err != nil {
		if err.Error() == "metric not found" {
			middleware.RespondWithError(w, http.StatusNotFound, "Metric not found", err)
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update metric", err)
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, updatedMetric)
}

// DeleteMetric deletes a metric by ID
// @Summary Delete a metric
// @Description Delete a metric by its ID
// @Tags metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Metric ID"
// @Success 204 "No content"
// @Failure 400 {object} middleware.ErrorResponse "Invalid ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 404 {object} middleware.ErrorResponse "Not found"
// @Failure 500 {object} middleware.ErrorResponse "Server error"
// @Router /metrics/{id} [delete]
func (h *MetricHandler) DeleteMetric(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	metricID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid metric ID", err)
		return
	}

	err = h.service.DeleteMetric(metricID)
	if err != nil {
		if err.Error() == "metric not found" {
			middleware.RespondWithError(w, http.StatusNotFound, "Metric not found", err)
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete metric", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
