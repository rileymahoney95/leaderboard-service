package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
