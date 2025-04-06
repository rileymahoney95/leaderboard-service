package models

import (
	"leaderboard-service/enums"
)

// Metric defines the types of measurable values used in leaderboards
type Metric struct {
	BaseModel
	Name            string                `gorm:"not null"`
	Description     string                `gorm:"type:text"`
	DataType        enums.MetricDataType  `gorm:"not null"` // e.g., "integer", "decimal", "boolean"
	Unit            string                // e.g., "calls", "texts", "%"
	AggregationType enums.AggregationType `gorm:"not null"` // e.g., "sum", "average", "count"
	ResetPeriod     enums.ResetPeriod     `gorm:"not null"` // e.g., "none", "daily", "weekly", "monthly", "yearly"
	IsHigherBetter  bool                  `gorm:"not null"`

	// Association to MetricValues
	Values []MetricValue `gorm:"foreignKey:MetricID;references:ID"`
}
