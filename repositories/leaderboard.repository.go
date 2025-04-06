package repositories

import (
	"leaderboard-service/db"
	"leaderboard-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaderboardRepository interface {
	Create(leaderboard *models.Leaderboard) error
	FindByID(id uuid.UUID) (*models.Leaderboard, error)
	FindAll() ([]models.Leaderboard, error)
	Update(leaderboard *models.Leaderboard) error
	Delete(id uuid.UUID) error
}

type leaderboardRepository struct {
	db *gorm.DB
}

func NewLeaderboardRepository() LeaderboardRepository {
	return &leaderboardRepository{
		db: db.DB,
	}
}

func (r *leaderboardRepository) Create(leaderboard *models.Leaderboard) error {
	return r.db.Create(leaderboard).Error
}

func (r *leaderboardRepository) FindByID(id uuid.UUID) (*models.Leaderboard, error) {
	var leaderboard models.Leaderboard
	err := r.db.First(&leaderboard, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

func (r *leaderboardRepository) FindAll() ([]models.Leaderboard, error) {
	var leaderboards []models.Leaderboard
	err := r.db.Find(&leaderboards).Error
	return leaderboards, err
}

func (r *leaderboardRepository) Update(leaderboard *models.Leaderboard) error {
	return r.db.Save(leaderboard).Error
}

func (r *leaderboardRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Leaderboard{}, "id = ?", id).Error
}
