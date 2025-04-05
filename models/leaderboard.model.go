package models

import (
	"time"
)

type Leaderboard struct {
	BaseModel
	Name            string `gorm:"not null"`
	Description     string `gorm:"type:text"`
	Category        string `gorm:"not null"`
	Type            string `gorm:"not null"` // individual, team
	TimeFrame       string `gorm:"not null"` // daily, weekly, monthly, yearly, all-time
	StartDate       *time.Time
	EndDate         *time.Time
	SortOrder       string `gorm:"not null"` // ascending, descending
	VisibilityScope string `gorm:"not null"` // public, private
	MaxEntries      int
	IsActive        bool

	Metrics []LeaderboardMetric `gorm:"foreignKey:LeaderboardID;references:ID"`
	Entries []LeaderboardEntry  `gorm:"foreignKey:LeaderboardID;references:ID"`
}
