package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"leaderboard-service/db"
	"leaderboard-service/enums"
	"leaderboard-service/models"
	"leaderboard-service/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CreateLeaderboardRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	Type            string  `json:"type,Type"`
	TimeFrame       string  `json:"time_frame,TimeFrame"`
	StartDate       *string `json:"start_date,StartDate"`
	EndDate         *string `json:"end_date,EndDate"`
	SortOrder       string  `json:"sort_order,SortOrder"`
	VisibilityScope string  `json:"visibility_scope,VisibilityScope"`
	IsActive        bool    `json:"is_active,IsActive"`
	MaxEntries      int     `json:"max_entries,MaxEntries"`
}

type UpdateLeaderboardRequest struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	Category        *string `json:"category"`
	Type            *string `json:"type,Type"`
	TimeFrame       *string `json:"time_frame,TimeFrame"`
	StartDate       *string `json:"start_date,StartDate"`
	EndDate         *string `json:"end_date,EndDate"`
	SortOrder       *string `json:"sort_order,SortOrder"`
	VisibilityScope *string `json:"visibility_scope,VisibilityScope"`
	IsActive        *bool   `json:"is_active,IsActive"`
	MaxEntries      *int    `json:"max_entries,MaxEntries"`
}

func CreateLeaderboard(w http.ResponseWriter, r *http.Request) {
	var req CreateLeaderboardRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := validateCreateLeaderboardRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert string to enum and validate
	leaderboardType := enums.LeaderboardType(req.Type)
	if !leaderboardType.Valid() {
		errorMessage := "Invalid leaderboard type, valid values are: " + strings.Join(enums.GetValidLeaderboardTypes(), ", ")
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	timeFrame := enums.TimeFrame(req.TimeFrame)
	if !timeFrame.Valid() {
		errorMessage := "Invalid leaderboard time frame, valid values are: " + strings.Join(enums.GetValidTimeFrames(), ", ")
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	sortOrder := enums.SortOrder(req.SortOrder)
	if !sortOrder.Valid() {
		errorMessage := "Invalid leaderboard sort order, valid values are: " + strings.Join(enums.GetValidSortOrders(), ", ")
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	visibilityScope := enums.VisibilityScope(req.VisibilityScope)
	if !visibilityScope.Valid() {
		errorMessage := "Invalid leaderboard visibility scope, valid values are: " + strings.Join(enums.GetValidVisibilityScopes(), ", ")
		http.Error(w, errorMessage, http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, leaderboard)
}

// validateCreateLeaderboardRequest validates required fields for creating a leaderboard
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

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	leaderboardId, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid leaderboard ID", http.StatusBadRequest)
		return
	}

	leaderboard := models.Leaderboard{}
	db.DB.First(&leaderboard, "id = ?", leaderboardId)

	if leaderboard.ID == uuid.Nil {
		http.Error(w, "Leaderboard not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, leaderboard)
}

func ListLeaderboards(w http.ResponseWriter, r *http.Request) {
	leaderboards := []models.Leaderboard{}
	db.DB.Find(&leaderboards)

	render.JSON(w, r, leaderboards)
}

func UpdateLeaderboard(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	leaderboardID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid leaderboard ID", http.StatusBadRequest)
		return
	}

	// Read the raw request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Fetch existing leaderboard
	var leaderboard models.Leaderboard
	if err := db.DB.First(&leaderboard, "id = ?", leaderboardID).Error; err != nil {
		http.Error(w, "Leaderboard not found", http.StatusNotFound)
		return
	}

	// Parse the JSON request
	var rawData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &rawData); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
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
			http.Error(w, "Invalid leaderboard type", http.StatusBadRequest)
			return
		}
		leaderboard.Type = leaderboardType
	}
	if req.TimeFrame != nil {
		timeFrame := enums.TimeFrame(*req.TimeFrame)
		if !timeFrame.Valid() {
			http.Error(w, "Invalid leaderboard time frame", http.StatusBadRequest)
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
			http.Error(w, "Invalid leaderboard sort order", http.StatusBadRequest)
			return
		}
		leaderboard.SortOrder = sortOrder
	}
	if req.VisibilityScope != nil {
		visibilityScope := enums.VisibilityScope(*req.VisibilityScope)
		if !visibilityScope.Valid() {
			http.Error(w, "Invalid leaderboard visibility scope", http.StatusBadRequest)
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
		http.Error(w, "Failed to update leaderboard", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, leaderboard)
}

// processRequestFields processes JSON fields with case-insensitive matching
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

func DeleteLeaderboard(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	leaderboardID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid leaderboard ID", http.StatusBadRequest)
		return
	}

	leaderboard := models.Leaderboard{}
	if err := db.DB.First(&leaderboard, "id = ?", leaderboardID).Error; err != nil {
		http.Error(w, "Leaderboard not found", http.StatusNotFound)
		return
	}

	if err := db.DB.Delete(&models.Leaderboard{}, "id = ?", leaderboardID).Error; err != nil {
		http.Error(w, "Failed to delete leaderboard", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
