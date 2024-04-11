package dto

import (
	"homework/internal/domain"
	"time"
)

// Sensor - структура для хранения данных датчика
type Sensor struct {
	ID           *int64            `validate:"min:0" json:"sensor_id"`
	SerialNumber *string           `validate:"min:10" json:"serial_number"`
	Type         domain.SensorType `validate:"in:SensorType" json:"type"`
	CurrentState int64             `validate:"min:0" json:"user_id"`
	Description  string            `validate:"min:0" json:"description"`
	IsActive     bool              `json:"is_active"`
	RegisteredAt time.Time         `json:"registered_at"`
	LastActivity time.Time         `json:"last_activity"`
}

func (s *Sensor) InitData() {
	s.initID()
	s.initRegisteredAt()
	s.initLastActivity()
	s.initType()
	s.initSerialNumber()
}

func (s *Sensor) initID() {
	if s.ID == nil {
		s.ID = new(int64)
		*s.ID = 1
	}
}

func (s *Sensor) initRegisteredAt() {
	if s.RegisteredAt.IsZero() {
		s.RegisteredAt = time.Now()
	}
}

func (s *Sensor) initLastActivity() {
	if s.LastActivity.IsZero() {
		s.LastActivity = time.Now()
	}
}

func (s *Sensor) initType() {
	if s.Type == "" {
		s.Type = domain.SensorTypeContactClosure
	}
}

func (s *Sensor) initSerialNumber() {
	if s.SerialNumber == nil {
		s.SerialNumber = new(string)
		*s.SerialNumber = "0123456789"
	}
}
