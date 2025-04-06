package repositories

import (
	"leaderboard-service/db"
	"leaderboard-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParticipantRepository interface {
	Create(participant *models.Participant) error
	FindByID(id uuid.UUID) (*models.Participant, error)
	FindAll() ([]models.Participant, error)
	Update(participant *models.Participant) error
	Delete(id uuid.UUID) error
}

type participantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository() ParticipantRepository {
	return &participantRepository{
		db: db.DB,
	}
}

func (r *participantRepository) Create(participant *models.Participant) error {
	return r.db.Create(participant).Error
}

func (r *participantRepository) FindByID(id uuid.UUID) (*models.Participant, error) {
	var participant models.Participant
	err := r.db.First(&participant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *participantRepository) FindAll() ([]models.Participant, error) {
	var participants []models.Participant
	err := r.db.Find(&participants).Error
	return participants, err
}

func (r *participantRepository) Update(participant *models.Participant) error {
	return r.db.Save(participant).Error
}

func (r *participantRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Participant{}, "id = ?", id).Error
}
