package models

type Participant struct {
	BaseModel
	ExternalID string      `gorm:"index"`
	Name       string      `gorm:"not null"`
	Type       string      `gorm:"not null"` // individual, team, group
	Metadata   interface{} `gorm:"type:jsonb"`

	// Association to MetricValues
	MetricValues []MetricValue `gorm:"foreignKey:ParticipantID;references:ID"`
}
