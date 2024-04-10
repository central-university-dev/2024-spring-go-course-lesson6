package inmemory

import (
	"context"
	"homework/internal/domain"
	"sync"
)

type SensorOwnerRepository struct {
	mu          sync.RWMutex
	SensorOwner []domain.SensorOwner
}

func NewSensorOwnerRepository() *SensorOwnerRepository {
	return &SensorOwnerRepository{SensorOwner: make([]domain.SensorOwner, 0)}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	r.mu.Lock()
	r.SensorOwner = append(r.SensorOwner, sensorOwner)
	r.mu.Unlock()
	return nil
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	sensorOwners := []domain.SensorOwner{}
	r.mu.RLock()
	for _, so := range r.SensorOwner {
		if so.UserID == userID {
			sensorOwners = append(sensorOwners, so)
		}
	}
	r.mu.RUnlock()
	return sensorOwners, nil
}
