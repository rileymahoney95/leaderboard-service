package services

import (
	"errors"
	"leaderboard-service/models"
	"leaderboard-service/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetricValueService interface {
	CreateMetricValue(metricID, participantID uuid.UUID, value float64, timestamp time.Time,
		source string, context interface{}) (*models.MetricValue, error)
	GetMetricValue(id uuid.UUID) (*models.MetricValue, error)
	ListMetricValues() ([]models.MetricValue, error)
	ListFilteredMetricValues(metricID, participantID *uuid.UUID, fromTime, toTime *time.Time) ([]models.MetricValue, error)
	UpdateMetricValue(id uuid.UUID, value *float64, timestamp *time.Time, source *string,
		context *interface{}) (*models.MetricValue, error)
	DeleteMetricValue(id uuid.UUID) error

	// Extra methods that verify entity existence
	VerifyMetricExists(metricID uuid.UUID) error
	VerifyParticipantExists(participantID uuid.UUID) error
}

type metricValueService struct {
	repo            repositories.MetricValueRepository
	metricRepo      repositories.MetricRepository
	participantRepo repositories.ParticipantRepository
}

func NewMetricValueService(repo repositories.MetricValueRepository,
	metricRepo repositories.MetricRepository,
	participantRepo repositories.ParticipantRepository) MetricValueService {
	return &metricValueService{
		repo:            repo,
		metricRepo:      metricRepo,
		participantRepo: participantRepo,
	}
}

func (s *metricValueService) CreateMetricValue(metricID, participantID uuid.UUID, value float64,
	timestamp time.Time, source string, context interface{}) (*models.MetricValue, error) {

	// Verify metric exists
	if err := s.VerifyMetricExists(metricID); err != nil {
		return nil, err
	}

	// Verify participant exists
	if err := s.VerifyParticipantExists(participantID); err != nil {
		return nil, err
	}

	// Set timestamp to current time if not provided
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	metricValue := models.MetricValue{
		MetricID:      metricID,
		ParticipantID: participantID,
		Value:         value,
		Timestamp:     timestamp,
		Source:        source,
		Context:       context,
	}

	err := s.repo.Create(&metricValue)
	if err != nil {
		return nil, err
	}

	return &metricValue, nil
}

func (s *metricValueService) GetMetricValue(id uuid.UUID) (*models.MetricValue, error) {
	metricValue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("metric value not found")
		}
		return nil, err
	}
	return metricValue, nil
}

func (s *metricValueService) ListMetricValues() ([]models.MetricValue, error) {
	return s.repo.FindAll()
}

func (s *metricValueService) ListFilteredMetricValues(metricID, participantID *uuid.UUID,
	fromTime, toTime *time.Time) ([]models.MetricValue, error) {
	return s.repo.FindFiltered(metricID, participantID, fromTime, toTime)
}

func (s *metricValueService) UpdateMetricValue(id uuid.UUID, value *float64, timestamp *time.Time,
	source *string, context *interface{}) (*models.MetricValue, error) {

	metricValue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("metric value not found")
		}
		return nil, err
	}

	// Apply the updates to the metric value
	if value != nil {
		metricValue.Value = *value
	}
	if timestamp != nil {
		metricValue.Timestamp = *timestamp
	}
	if source != nil {
		metricValue.Source = *source
	}
	if context != nil {
		metricValue.Context = *context
	}

	err = s.repo.Update(metricValue)
	if err != nil {
		return nil, err
	}

	return metricValue, nil
}

func (s *metricValueService) DeleteMetricValue(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("metric value not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

// Verify that a metric exists
func (s *metricValueService) VerifyMetricExists(metricID uuid.UUID) error {
	_, err := s.metricRepo.FindByID(metricID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("metric not found")
		}
		return err
	}
	return nil
}

// Verify that a participant exists
func (s *metricValueService) VerifyParticipantExists(participantID uuid.UUID) error {
	_, err := s.participantRepo.FindByID(participantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("participant not found")
		}
		return err
	}
	return nil
}
