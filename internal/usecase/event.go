package usecase

import (
	"context"
	"homework/internal/domain"
	"time"
)

type Event struct {
	EventRepo  EventRepository
	SensorRepo SensorRepository
}

func NewEvent(er EventRepository, sr SensorRepository) *Event {
	return &Event{EventRepo: er, SensorRepo: sr}
}

func validateEvent(e *domain.Event) error {
	if e.Timestamp.IsZero() {
		return ErrInvalidEventTimestamp
	}
	return nil
}

func (e *Event) ReceiveEvent(ctx context.Context, event *domain.Event) error {
	if err := validateEvent(event); err != nil {
		return err
	}

	sensor, err := e.SensorRepo.GetSensorBySerialNumber(ctx, event.SensorSerialNumber)
	if err != nil {
		return err
	}
	sensor.LastActivity = time.Now()
	sensor.CurrentState = event.Payload
	event.SensorID = sensor.ID

	if err := e.EventRepo.SaveEvent(ctx, event); err != nil {
		return err
	}
	return e.SensorRepo.SaveSensor(ctx, sensor)
}

func (e *Event) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	if _, err := e.SensorRepo.GetSensorByID(ctx, id); err != nil {
		return nil, err
	}
	return e.EventRepo.GetLastEventBySensorID(ctx, id)
}
