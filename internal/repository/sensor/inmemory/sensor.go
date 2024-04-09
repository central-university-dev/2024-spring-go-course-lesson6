package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrSensorNil      = errors.New("sensor is nil")
	ErrSensorNotFound = errors.New("sensor not found")
)

type SensorRepository struct {
	mu            sync.RWMutex
	SensorStorage map[int64]domain.Sensor
}

func NewSensorRepository() *SensorRepository {
	return &SensorRepository{SensorStorage: make(map[int64]domain.Sensor)}
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	now := time.Now()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if sensor == nil {
		return ErrSensorNil
	}

	sensor.RegisteredAt = now
	sensor.ID = rand.Int63()
	r.mu.Lock()
	r.SensorStorage[sensor.ID] = *sensor
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
	if s, ok := (r.SensorStorage)[id]; ok {
		return &s, nil
	}
	r.mu.RUnlock()
	return nil, ErrSensorNotFound
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.mu.RLock()
	for _, s := range r.SensorStorage {
		if s.SerialNumber == sn {
			return &s, nil
		}
	}
	r.mu.RUnlock()
	return nil, ErrSensorNotFound
}
