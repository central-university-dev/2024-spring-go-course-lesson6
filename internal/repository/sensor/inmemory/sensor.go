package inmemory

import (
	"context"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
	"time"
)

type SensorRepository struct {
	mu            sync.RWMutex
	SensorStorage []domain.Sensor
}

func NewSensorRepository() *SensorRepository {
	return &SensorRepository{SensorStorage: make([]domain.Sensor, 0)}
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if sensor == nil {
		return usecase.ErrSensorNotFound
	}

	sensor.RegisteredAt = time.Now()
	sensor.ID = 1

	r.mu.Lock()
	r.SensorStorage = append(r.SensorStorage, *sensor)
	r.mu.Unlock()

	return nil
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	var sensors []domain.Sensor
	if ctx.Err() != nil {
		return sensors, ctx.Err()
	}

	sensors = make([]domain.Sensor, len(r.SensorStorage))
	i := 0
	r.mu.RLock()
	for _, s := range r.SensorStorage {
		sensors[i] = s
		i++
	}
	r.mu.RUnlock()
	return sensors, nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.SensorStorage {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, usecase.ErrSensorNotFound
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.SensorStorage {
		if s.SerialNumber == sn {
			return &s, nil
		}
	}
	return nil, usecase.ErrSensorNotFound
}
