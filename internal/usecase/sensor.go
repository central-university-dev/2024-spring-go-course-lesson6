package usecase

import (
	"context"
	"errors"
	"homework/internal/domain"
)

type Sensor struct {
	SensorRepo SensorRepository
}

func NewSensor(sr SensorRepository) *Sensor {
	return &Sensor{SensorRepo: sr}
}

func validateSensor(s *domain.Sensor) error {
	if s.Type != domain.SensorTypeContactClosure && s.Type != domain.SensorTypeADC {
		return ErrWrongSensorType
	} else if len(s.SerialNumber) != 10 {
		return ErrWrongSensorSerialNumber
	}
	return nil
}

func (s *Sensor) RegisterSensor(ctx context.Context, sensor *domain.Sensor) (*domain.Sensor, error) {
	if err := validateSensor(sensor); err != nil {
		return nil, err
	}

	sensorDb, err := s.SensorRepo.GetSensorBySerialNumber(ctx, sensor.SerialNumber)
	if err == nil {
		return sensorDb, nil
	} else if errors.Is(err, ErrSensorNotFound) {
		if err := s.SensorRepo.SaveSensor(ctx, sensor); err != nil {
			return nil, err
		}
		return sensor, nil
	}
	return nil, err
}

func (s *Sensor) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	sensors, err := s.SensorRepo.GetSensors(ctx)
	if err != nil {
		return nil, err
	}
	return sensors, nil
}

func (s *Sensor) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	sensors, err := s.SensorRepo.GetSensorByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return sensors, nil
}
