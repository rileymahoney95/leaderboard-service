package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"leaderboard-service/db"
	"leaderboard-service/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CreateLeaderboardRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	Type            string  `json:"type"`
	TimeFrame       string  `json:"time_frame"`
	StartDate       *string `json:"start_date"`
	EndDate         *string `json:"end_date"`
	SortOrder       string  `json:"sort_order"`
	VisibilityScope string  `json:"visibility_scope"`
	MaxEntries      int     `json:"max_entries"`
}

func CreateLeaderboard(w http.ResponseWriter, r *http.Request) {
	var req CreateLeaderboardRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse optional dates
	var startDate *time.Time
	var endDate *time.Time

	if req.StartDate != nil {
		parsedDate, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		startDate = &parsedDate
	}

	if req.EndDate != nil {
		parsedDate, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		endDate = &parsedDate
	}

	leaderboard := models.Leaderboard{
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		Type:            req.Type,
		TimeFrame:       req.TimeFrame,
		StartDate:       startDate,
		EndDate:         endDate,
		SortOrder:       req.SortOrder,
		VisibilityScope: req.VisibilityScope,
		MaxEntries:      req.MaxEntries,
	}

	err = db.DB.Create(&leaderboard).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, leaderboard)
}

// /{leaderboardId}
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
