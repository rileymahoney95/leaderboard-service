package services

import (
	"errors"
	"leaderboard-service/models"
	"leaderboard-service/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParticipantService interface {
	CreateParticipant(externalID, name, participantType string, metadata interface{}) (*models.Participant, error)
	GetParticipant(id uuid.UUID) (*models.Participant, error)
	ListParticipants() ([]models.Participant, error)
	UpdateParticipant(id uuid.UUID, externalID, name, participantType *string, metadata *interface{}) (*models.Participant, error)
	DeleteParticipant(id uuid.UUID) error
}

type participantService struct {
	repo repositories.ParticipantRepository
}

func NewParticipantService(repo repositories.ParticipantRepository) ParticipantService {
	return &participantService{
		repo: repo,
	}
}

func (s *participantService) CreateParticipant(externalID, name, participantType string, metadata interface{}) (*models.Participant, error) {
	participant := models.Participant{
		ExternalID: externalID,
		Name:       name,
		Type:       participantType,
		Metadata:   metadata,
	}

	err := s.repo.Create(&participant)
	if err != nil {
		return nil, err
	}

	return &participant, nil
}

func (s *participantService) GetParticipant(id uuid.UUID) (*models.Participant, error) {
	participant, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("participant not found")
		}
		return nil, err
	}
	return participant, nil
}

func (s *participantService) ListParticipants() ([]models.Participant, error) {
	return s.repo.FindAll()
}

func (s *participantService) UpdateParticipant(id uuid.UUID, externalID, name, participantType *string, metadata *interface{}) (*models.Participant, error) {
	participant, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("participant not found")
		}
		return nil, err
	}

	// Apply the updates to the participant
	if externalID != nil {
		participant.ExternalID = *externalID
	}
	if name != nil {
		participant.Name = *name
	}
	if participantType != nil {
		participant.Type = *participantType
	}
	if metadata != nil {
		participant.Metadata = *metadata
	}

	err = s.repo.Update(participant)
	if err != nil {
		return nil, err
	}

	return participant, nil
}

func (s *participantService) DeleteParticipant(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("participant not found")
		}
		return err
	}

	return s.repo.Delete(id)
}
