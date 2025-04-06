package repositories

import (
	"leaderboard-service/db"
	"leaderboard-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaderboardEntryRepository interface {
	Create(entry *models.LeaderboardEntry) error
	FindByID(id uuid.UUID) (*models.LeaderboardEntry, error)
	FindAll() ([]models.LeaderboardEntry, error)
	FindByLeaderboardID(leaderboardID uuid.UUID) ([]models.LeaderboardEntry, error)
	FindByParticipantID(participantID uuid.UUID) ([]models.LeaderboardEntry, error)
	FindFiltered(leaderboardID, participantID *uuid.UUID) ([]models.LeaderboardEntry, error)
	Update(entry *models.LeaderboardEntry) error
	Delete(id uuid.UUID) error
}

type leaderboardEntryRepository struct {
	db *gorm.DB
}

func NewLeaderboardEntryRepository() LeaderboardEntryRepository {
	return &leaderboardEntryRepository{
		db: db.DB,
	}
}

func (r *leaderboardEntryRepository) Create(entry *models.LeaderboardEntry) error {
	return r.db.Create(entry).Error
}

func (r *leaderboardEntryRepository) FindByID(id uuid.UUID) (*models.LeaderboardEntry, error) {
	var entry models.LeaderboardEntry
	err := r.db.First(&entry, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *leaderboardEntryRepository) FindAll() ([]models.LeaderboardEntry, error) {
	var entries []models.LeaderboardEntry
	err := r.db.Find(&entries).Error
	return entries, err
}

func (r *leaderboardEntryRepository) FindByLeaderboardID(leaderboardID uuid.UUID) ([]models.LeaderboardEntry, error) {
	var entries []models.LeaderboardEntry
	err := r.db.Where("leaderboard_id = ?", leaderboardID).Order("rank asc").Find(&entries).Error
	return entries, err
}

func (r *leaderboardEntryRepository) FindByParticipantID(participantID uuid.UUID) ([]models.LeaderboardEntry, error) {
	var entries []models.LeaderboardEntry
	err := r.db.Where("participant_id = ?", participantID).Find(&entries).Error
	return entries, err
}

func (r *leaderboardEntryRepository) FindFiltered(leaderboardID, participantID *uuid.UUID) ([]models.LeaderboardEntry, error) {
	var entries []models.LeaderboardEntry
	query := r.db

	if leaderboardID != nil {
		query = query.Where("leaderboard_id = ?", *leaderboardID)
	}

	if participantID != nil {
		query = query.Where("participant_id = ?", *participantID)
	}

	// Order by rank
	err := query.Order("rank asc").Find(&entries).Error
	return entries, err
}

func (r *leaderboardEntryRepository) Update(entry *models.LeaderboardEntry) error {
	return r.db.Save(entry).Error
}

func (r *leaderboardEntryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.LeaderboardEntry{}, "id = ?", id).Error
}
