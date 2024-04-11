package dto

import (
	"time"
)

// Event - структура события по датчику
type Event struct {
	Timestamp          time.Time `json:"timestamp"`
	SensorSerialNumber string    `validate:"min:10" json:"sensor_serial_number"`
	SensorID           int64     `validate:"min:0" json:"sensor_id"`
	Payload            int64     `validate:"min:0" json:"payload"`
}

func (e *Event) InitData() {
	e.initTimestamp()
}

func (e *Event) initTimestamp() {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
}
