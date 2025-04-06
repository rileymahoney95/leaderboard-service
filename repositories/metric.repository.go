package repositories

import (
	"leaderboard-service/db"
	"leaderboard-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetricRepository interface {
	Create(metric *models.Metric) error
	FindByID(id uuid.UUID) (*models.Metric, error)
	FindAll() ([]models.Metric, error)
	Update(metric *models.Metric) error
	Delete(id uuid.UUID) error
}

type metricRepository struct {
	db *gorm.DB
}

func NewMetricRepository() MetricRepository {
	return &metricRepository{
		db: db.DB,
	}
}

func (r *metricRepository) Create(metric *models.Metric) error {
	return r.db.Create(metric).Error
}

func (r *metricRepository) FindByID(id uuid.UUID) (*models.Metric, error) {
	var metric models.Metric
	err := r.db.First(&metric, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *metricRepository) FindAll() ([]models.Metric, error) {
	var metrics []models.Metric
	err := r.db.Find(&metrics).Error
	return metrics, err
}

func (r *metricRepository) Update(metric *models.Metric) error {
	return r.db.Save(metric).Error
}

func (r *metricRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Metric{}, "id = ?", id).Error
}
