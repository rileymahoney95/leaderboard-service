package models

import (
	"github.com/google/uuid"
)

// LeaderboardMetric represents a metric associated with a leaderboard
type LeaderboardMetric struct {
	BaseModel
	LeaderboardID   uuid.UUID `gorm:"type:uuid;not null"`
	MetricID        uuid.UUID `gorm:"type:uuid;not null"`
	Weight          float64   `gorm:"not null;default:1.0"`
	DisplayPriority int       `gorm:"not null;default:0"`
}
