package domain

import (
	"errors"
	"time"
)

var (
	ErrWrongSensorType         = errors.New("wrong sensor type")
	ErrWrongSensorSerialNumber = errors.New("wrong sensor serial number")
)

type SensorType string

const (
	SensorTypeContactClosure SensorType = "cc"
	SensorTypeADC            SensorType = "adc"
)

// Sensor - структура для хранения данных датчика
type Sensor struct {
	ID           int64
	SerialNumber string
	Type         SensorType
	CurrentState int64
	Description  string
	IsActive     bool
	RegisteredAt time.Time
	LastActivity time.Time
}

func (s *Sensor) Validate() error {
	if s.Type != SensorTypeContactClosure && s.Type != SensorTypeADC {
		return ErrWrongSensorType
	} else if len(s.SerialNumber) != 10 {
		return ErrWrongSensorSerialNumber
	}
	return nil
}
