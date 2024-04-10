package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"sync"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventNil      = errors.New("event is nil")
)

type EventRepository struct {
	mu           sync.RWMutex
	EventStorage map[int64][]domain.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{EventStorage: make(map[int64][]domain.Event)}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if event == nil {
		return ErrEventNil
	}

	r.mu.Lock()
	if _, ok := r.EventStorage[event.SensorID]; !ok {
		r.EventStorage[event.SensorID] = []domain.Event{*event}
	} else {
		r.EventStorage[event.SensorID] = append(r.EventStorage[event.SensorID], *event)
	}
	r.mu.Unlock()
	return nil
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if eventStorages, ok := r.EventStorage[id]; ok {
		lstEvent := eventStorages[0]
		for _, e := range eventStorages {
			if e.Timestamp.After(lstEvent.Timestamp) {
				lstEvent = e
			}
		}
		return &lstEvent, nil
	}
	return nil, ErrEventNotFound
}
