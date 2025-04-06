package repositories

import (
	"leaderboard-service/db"
	"leaderboard-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetricValueRepository interface {
	Create(metricValue *models.MetricValue) error
	FindByID(id uuid.UUID) (*models.MetricValue, error)
	FindAll() ([]models.MetricValue, error)
	FindByMetricID(metricID uuid.UUID) ([]models.MetricValue, error)
	FindByParticipantID(participantID uuid.UUID) ([]models.MetricValue, error)
	FindFiltered(metricID, participantID *uuid.UUID, fromTime, toTime *time.Time) ([]models.MetricValue, error)
	Update(metricValue *models.MetricValue) error
	Delete(id uuid.UUID) error
}

type metricValueRepository struct {
	db *gorm.DB
}

func NewMetricValueRepository() MetricValueRepository {
	return &metricValueRepository{
		db: db.DB,
	}
}

func (r *metricValueRepository) Create(metricValue *models.MetricValue) error {
	return r.db.Create(metricValue).Error
}

func (r *metricValueRepository) FindByID(id uuid.UUID) (*models.MetricValue, error) {
	var metricValue models.MetricValue
	err := r.db.First(&metricValue, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &metricValue, nil
}

func (r *metricValueRepository) FindAll() ([]models.MetricValue, error) {
	var metricValues []models.MetricValue
	err := r.db.Find(&metricValues).Error
	return metricValues, err
}

func (r *metricValueRepository) FindByMetricID(metricID uuid.UUID) ([]models.MetricValue, error) {
	var metricValues []models.MetricValue
	err := r.db.Where("metric_id = ?", metricID).Find(&metricValues).Error
	return metricValues, err
}

func (r *metricValueRepository) FindByParticipantID(participantID uuid.UUID) ([]models.MetricValue, error) {
	var metricValues []models.MetricValue
	err := r.db.Where("participant_id = ?", participantID).Find(&metricValues).Error
	return metricValues, err
}

func (r *metricValueRepository) FindFiltered(metricID, participantID *uuid.UUID, fromTime, toTime *time.Time) ([]models.MetricValue, error) {
	var metricValues []models.MetricValue
	query := r.db

	if metricID != nil {
		query = query.Where("metric_id = ?", *metricID)
	}

	if participantID != nil {
		query = query.Where("participant_id = ?", *participantID)
	}

	if fromTime != nil {
		query = query.Where("timestamp >= ?", *fromTime)
	}

	if toTime != nil {
		query = query.Where("timestamp <= ?", *toTime)
	}

	// Order by timestamp, most recent first
	err := query.Order("timestamp desc").Find(&metricValues).Error
	return metricValues, err
}

func (r *metricValueRepository) Update(metricValue *models.MetricValue) error {
	return r.db.Save(metricValue).Error
}

func (r *metricValueRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.MetricValue{}, "id = ?", id).Error
}
