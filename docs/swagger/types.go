package swagger

import (
	"time"

	"github.com/google/uuid"
)

// This file is only used for Swagger documentation
// It provides type definitions for Swagger to understand GORM types

// DeletedAt is used to represent GORM's DeletedAt field for soft deletes
type DeletedAt struct {
	Time  time.Time
	Valid bool
}

// BaseModel represents the base model fields used in all models
type BaseModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt DeletedAt `json:"deleted_at,omitempty"`
}

// LeaderboardType represents the type of a leaderboard
type LeaderboardType string

// TimeFrame represents the time frame of a leaderboard
type TimeFrame string

// SortOrder represents the sort order of a leaderboard
type SortOrder string

// VisibilityScope represents the visibility scope of a leaderboard
type VisibilityScope string

// Leaderboard represents a leaderboard entity for Swagger documentation
type Leaderboard struct {
	BaseModel
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Category        string          `json:"category"`
	Type            LeaderboardType `json:"type"`
	TimeFrame       TimeFrame       `json:"time_frame"`
	StartDate       *time.Time      `json:"start_date,omitempty"`
	EndDate         *time.Time      `json:"end_date,omitempty"`
	SortOrder       SortOrder       `json:"sort_order"`
	VisibilityScope VisibilityScope `json:"visibility_scope"`
	MaxEntries      int             `json:"max_entries"`
	IsActive        bool            `json:"is_active"`
}
