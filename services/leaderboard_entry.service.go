package services

import (
	"errors"
	"leaderboard-service/models"
	"leaderboard-service/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaderboardEntryService interface {
	CreateLeaderboardEntry(leaderboardID, participantID uuid.UUID, score float64, rank int, lastUpdated time.Time) (*models.LeaderboardEntry, error)
	GetLeaderboardEntry(id uuid.UUID) (*models.LeaderboardEntry, error)
	ListLeaderboardEntries() ([]models.LeaderboardEntry, error)
	ListFilteredLeaderboardEntries(leaderboardID, participantID *uuid.UUID) ([]models.LeaderboardEntry, error)
	UpdateLeaderboardEntry(id uuid.UUID, score *float64, rank *int, lastUpdated *time.Time) (*models.LeaderboardEntry, error)
	DeleteLeaderboardEntry(id uuid.UUID) error

	// Verification methods
	VerifyLeaderboardExists(leaderboardID uuid.UUID) error
	VerifyParticipantExists(participantID uuid.UUID) error
}

type leaderboardEntryService struct {
	repo            repositories.LeaderboardEntryRepository
	leaderboardRepo repositories.LeaderboardRepository
	participantRepo repositories.ParticipantRepository
}

func NewLeaderboardEntryService(repo repositories.LeaderboardEntryRepository,
	leaderboardRepo repositories.LeaderboardRepository,
	participantRepo repositories.ParticipantRepository) LeaderboardEntryService {
	return &leaderboardEntryService{
		repo:            repo,
		leaderboardRepo: leaderboardRepo,
		participantRepo: participantRepo,
	}
}

func (s *leaderboardEntryService) CreateLeaderboardEntry(leaderboardID, participantID uuid.UUID,
	score float64, rank int, lastUpdated time.Time) (*models.LeaderboardEntry, error) {

	// Verify leaderboard exists
	if err := s.VerifyLeaderboardExists(leaderboardID); err != nil {
		return nil, err
	}

	// Verify participant exists
	if err := s.VerifyParticipantExists(participantID); err != nil {
		return nil, err
	}

	// Set lastUpdated to current time if not provided
	if lastUpdated.IsZero() {
		lastUpdated = time.Now()
	}

	entry := models.LeaderboardEntry{
		LeaderboardID: leaderboardID,
		ParticipantID: participantID,
		Score:         score,
		Rank:          rank,
		LastUpdated:   lastUpdated,
	}

	err := s.repo.Create(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (s *leaderboardEntryService) GetLeaderboardEntry(id uuid.UUID) (*models.LeaderboardEntry, error) {
	entry, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("leaderboard entry not found")
		}
		return nil, err
	}
	return entry, nil
}

func (s *leaderboardEntryService) ListLeaderboardEntries() ([]models.LeaderboardEntry, error) {
	return s.repo.FindAll()
}

func (s *leaderboardEntryService) ListFilteredLeaderboardEntries(leaderboardID, participantID *uuid.UUID) ([]models.LeaderboardEntry, error) {
	return s.repo.FindFiltered(leaderboardID, participantID)
}

func (s *leaderboardEntryService) UpdateLeaderboardEntry(id uuid.UUID, score *float64,
	rank *int, lastUpdated *time.Time) (*models.LeaderboardEntry, error) {

	entry, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("leaderboard entry not found")
		}
		return nil, err
	}

	// Apply the updates to the entry
	if score != nil {
		entry.Score = *score
	}
	if rank != nil {
		entry.Rank = *rank
	}
	if lastUpdated != nil {
		entry.LastUpdated = *lastUpdated
	} else {
		// Update the LastUpdated field if not explicitly provided
		entry.LastUpdated = time.Now()
	}

	err = s.repo.Update(entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *leaderboardEntryService) DeleteLeaderboardEntry(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("leaderboard entry not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

// Verify that a leaderboard exists
func (s *leaderboardEntryService) VerifyLeaderboardExists(leaderboardID uuid.UUID) error {
	_, err := s.leaderboardRepo.FindByID(leaderboardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("leaderboard not found")
		}
		return err
	}
	return nil
}

// Verify that a participant exists
func (s *leaderboardEntryService) VerifyParticipantExists(participantID uuid.UUID) error {
	_, err := s.participantRepo.FindByID(participantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("participant not found")
		}
		return err
	}
	return nil
}
