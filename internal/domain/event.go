package domain

import (
	"errors"
	"time"
)

var ErrInvalidEventTimestamp = errors.New("invalid event timestamp")

// Event - структура события по датчику
type Event struct {
	Timestamp          time.Time
	SensorSerialNumber string
	SensorID           int64
	Payload            int64
}

func (e *Event) Validate() error {
	if e.Timestamp.IsZero() {
		return ErrInvalidEventTimestamp
	}
	return nil
}