package services

import (
	"errors"
	"leaderboard-service/enums"
	"leaderboard-service/models"
	"leaderboard-service/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetricService interface {
	CreateMetric(name, description string, dataType enums.MetricDataType, unit string,
		aggregationType enums.AggregationType, resetPeriod enums.ResetPeriod, isHigherBetter bool) (*models.Metric, error)
	GetMetric(id uuid.UUID) (*models.Metric, error)
	ListMetrics() ([]models.Metric, error)
	UpdateMetric(id uuid.UUID, name, description *string, dataType *enums.MetricDataType,
		unit *string, aggregationType *enums.AggregationType, resetPeriod *enums.ResetPeriod,
		isHigherBetter *bool) (*models.Metric, error)
	DeleteMetric(id uuid.UUID) error
}

type metricService struct {
	repo repositories.MetricRepository
}

func NewMetricService(repo repositories.MetricRepository) MetricService {
	return &metricService{
		repo: repo,
	}
}

func (s *metricService) CreateMetric(name, description string, dataType enums.MetricDataType, unit string,
	aggregationType enums.AggregationType, resetPeriod enums.ResetPeriod, isHigherBetter bool) (*models.Metric, error) {

	metric := models.Metric{
		Name:            name,
		Description:     description,
		DataType:        dataType,
		Unit:            unit,
		AggregationType: aggregationType,
		ResetPeriod:     resetPeriod,
		IsHigherBetter:  isHigherBetter,
	}

	err := s.repo.Create(&metric)
	if err != nil {
		return nil, err
	}

	return &metric, nil
}

func (s *metricService) GetMetric(id uuid.UUID) (*models.Metric, error) {
	metric, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("metric not found")
		}
		return nil, err
	}
	return metric, nil
}

func (s *metricService) ListMetrics() ([]models.Metric, error) {
	return s.repo.FindAll()
}

func (s *metricService) UpdateMetric(id uuid.UUID, name, description *string, dataType *enums.MetricDataType,
	unit *string, aggregationType *enums.AggregationType, resetPeriod *enums.ResetPeriod,
	isHigherBetter *bool) (*models.Metric, error) {

	metric, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("metric not found")
		}
		return nil, err
	}

	// Apply the updates to the metric
	if name != nil {
		metric.Name = *name
	}
	if description != nil {
		metric.Description = *description
	}
	if dataType != nil {
		metric.DataType = *dataType
	}
	if unit != nil {
		metric.Unit = *unit
	}
	if aggregationType != nil {
		metric.AggregationType = *aggregationType
	}
	if resetPeriod != nil {
		metric.ResetPeriod = *resetPeriod
	}
	if isHigherBetter != nil {
		metric.IsHigherBetter = *isHigherBetter
	}

	err = s.repo.Update(metric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *metricService) DeleteMetric(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("metric not found")
		}
		return err
	}

	return s.repo.Delete(id)
}
