package models

import (
	"time"

	"github.com/google/uuid"
)

// LeaderboardEntry represents an entry/ranking in a leaderboard
type LeaderboardEntry struct {
	BaseModel
	LeaderboardID uuid.UUID `gorm:"type:uuid;not null"`
	ParticipantID uuid.UUID `gorm:"type:uuid;not null"`
	Rank          int       `gorm:"not null"`
	Score         float64   `gorm:"not null"`
	LastUpdated   time.Time `gorm:"not null"`
}
