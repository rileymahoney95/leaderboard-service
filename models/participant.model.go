package models

type Participant struct {
	BaseModel
	ExternalID string      `gorm:"index"`
	Name       string      `gorm:"not null"`
	Type       string      `gorm:"not null"` // individual, team, group
	Metadata   interface{} `gorm:"type:jsonb"`
}
