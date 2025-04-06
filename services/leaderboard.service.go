package services

import (
	"errors"
	"leaderboard-service/enums"
	"leaderboard-service/models"
	"leaderboard-service/repositories"
	"leaderboard-service/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaderboardService interface {
	CreateLeaderboard(name, description, category string, leaderboardType enums.LeaderboardType,
		timeFrame enums.TimeFrame, startDate, endDate *string, sortOrder enums.SortOrder,
		visibilityScope enums.VisibilityScope, maxEntries int, isActive bool) (*models.Leaderboard, error)
	GetLeaderboard(id uuid.UUID) (*models.Leaderboard, error)
	ListLeaderboards() ([]models.Leaderboard, error)
	UpdateLeaderboard(id uuid.UUID, name, description, category *string, leaderboardType *enums.LeaderboardType,
		timeFrame *enums.TimeFrame, startDate, endDate *string, sortOrder *enums.SortOrder,
		visibilityScope *enums.VisibilityScope, maxEntries *int, isActive *bool) (*models.Leaderboard, error)
	DeleteLeaderboard(id uuid.UUID) error
}

type leaderboardService struct {
	repo repositories.LeaderboardRepository
}

func NewLeaderboardService(repo repositories.LeaderboardRepository) LeaderboardService {
	return &leaderboardService{
		repo: repo,
	}
}

func (s *leaderboardService) CreateLeaderboard(name, description, category string, leaderboardType enums.LeaderboardType,
	timeFrame enums.TimeFrame, startDate, endDate *string, sortOrder enums.SortOrder,
	visibilityScope enums.VisibilityScope, maxEntries int, isActive bool) (*models.Leaderboard, error) {

	start, end := utils.ValidateDates(startDate, endDate)

	leaderboard := models.Leaderboard{
		Name:            name,
		Description:     description,
		Category:        category,
		Type:            leaderboardType,
		TimeFrame:       timeFrame,
		StartDate:       start,
		EndDate:         end,
		SortOrder:       sortOrder,
		VisibilityScope: visibilityScope,
		MaxEntries:      maxEntries,
		IsActive:        isActive,
	}

	err := s.repo.Create(&leaderboard)
	if err != nil {
		return nil, err
	}

	return &leaderboard, nil
}

func (s *leaderboardService) GetLeaderboard(id uuid.UUID) (*models.Leaderboard, error) {
	leaderboard, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("leaderboard not found")
		}
		return nil, err
	}
	return leaderboard, nil
}

func (s *leaderboardService) ListLeaderboards() ([]models.Leaderboard, error) {
	return s.repo.FindAll()
}

func (s *leaderboardService) UpdateLeaderboard(id uuid.UUID, name, description, category *string,
	leaderboardType *enums.LeaderboardType, timeFrame *enums.TimeFrame,
	startDate, endDate *string, sortOrder *enums.SortOrder,
	visibilityScope *enums.VisibilityScope, maxEntries *int, isActive *bool) (*models.Leaderboard, error) {

	leaderboard, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("leaderboard not found")
		}
		return nil, err
	}

	// Apply the updates to the leaderboard
	if name != nil {
		leaderboard.Name = *name
	}
	if description != nil {
		leaderboard.Description = *description
	}
	if category != nil {
		leaderboard.Category = *category
	}
	if leaderboardType != nil {
		leaderboard.Type = *leaderboardType
	}
	if timeFrame != nil {
		leaderboard.TimeFrame = *timeFrame
	}
	if startDate != nil || endDate != nil {
		start, end := utils.ValidateDates(startDate, endDate)
		if startDate != nil {
			leaderboard.StartDate = start
		}
		if endDate != nil {
			leaderboard.EndDate = end
		}
	}
	if sortOrder != nil {
		leaderboard.SortOrder = *sortOrder
	}
	if visibilityScope != nil {
		leaderboard.VisibilityScope = *visibilityScope
	}
	if maxEntries != nil {
		leaderboard.MaxEntries = *maxEntries
	}
	if isActive != nil {
		leaderboard.IsActive = *isActive
	}

	err = s.repo.Update(leaderboard)
	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}

func (s *leaderboardService) DeleteLeaderboard(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("leaderboard not found")
		}
		return err
	}

	return s.repo.Delete(id)
}
