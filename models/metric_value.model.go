package models

import (
	"time"

	"github.com/google/uuid"
)

// MetricValue stores actual recorded values for metrics for each participant
type MetricValue struct {
	BaseModel
	MetricID      uuid.UUID   `gorm:"type:uuid;not null"`
	ParticipantID uuid.UUID   `gorm:"type:uuid;not null"`
	Value         float64     `gorm:"not null"`
	Timestamp     time.Time   `gorm:"not null"`
	Source        string      // Identifies where/how this value was recorded
	Context       interface{} `gorm:"type:jsonb"` // For any additional data (e.g., distinguishing call vs. text)

	// Relations
	Metric      Metric      `gorm:"foreignKey:MetricID"`
	Participant Participant `gorm:"foreignKey:ParticipantID"`
}
