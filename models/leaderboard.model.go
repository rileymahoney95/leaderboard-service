package models

import (
	"leaderboard-service/enums"
	"time"
)

type Leaderboard struct {
	BaseModel
	Name            string                `gorm:"not null"`
	Description     string                `gorm:"type:text"`
	Category        string                `gorm:"not null"`
	Type            enums.LeaderboardType `gorm:"not null"`
	TimeFrame       enums.TimeFrame       `gorm:"not null"`
	StartDate       *time.Time
	EndDate         *time.Time
	SortOrder       enums.SortOrder       `gorm:"not null"`
	VisibilityScope enums.VisibilityScope `gorm:"not null"`
	MaxEntries      int
	IsActive        bool

	Metrics []LeaderboardMetric `gorm:"foreignKey:LeaderboardID;references:ID"`
	Entries []LeaderboardEntry  `gorm:"foreignKey:LeaderboardID;references:ID"`
}
