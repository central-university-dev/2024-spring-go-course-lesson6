package dto

// User - структура для хранения пользователя
type User struct {
	ID   int64  `json:"id"`
	Name string `validate:"min:0" json:"name"`
}

// SensorOwner - структура для связи пользователя и датчика
// UserID - id пользователя
// SensorID - id датчика
// Связь многие-ко-многим: пользователь может иметь доступ к нескольким датчикам, датчик может быть доступен для нескольких пользователей.
type SensorOwner struct {
	UserID   int64 `validate:"min:0" json:"user_id"`
	SensorID int64 `validate:"min:0" json:"sensor_id"`
}

func (u *User) InitData() {
	u.initID()
}

func (u *User) initID() {
	if u.ID == 0 {
		u.ID = 1
	}
}
