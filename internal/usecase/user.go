package usecase

import (
	"context"
	"homework/internal/domain"
)

type User struct {
	UserRepo        UserRepository
	SensorOwnerRepo SensorOwnerRepository
	SensorRepo      SensorRepository
}

func NewUser(ur UserRepository, sor SensorOwnerRepository, sr SensorRepository) *User {
	return &User{UserRepo: ur, SensorOwnerRepo: sor, SensorRepo: sr}
}

func validateUser(u *domain.User) error {
	if u.Name == "" {
		return ErrInvalidUserName
	}
	return nil
}

func (u *User) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := validateUser(user); err != nil {
		return nil, err
	}

	if err := u.UserRepo.SaveUser(ctx, user); err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) AttachSensorToUser(ctx context.Context, userID, sensorID int64) error {
	if _, err := u.UserRepo.GetUserByID(ctx, userID); err != nil {
		return err
	}

	sensor, err := u.SensorRepo.GetSensorByID(ctx, sensorID)
	if err != nil {
		return err
	}

	return u.SensorOwnerRepo.SaveSensorOwner(ctx, domain.SensorOwner{UserID: userID, SensorID: sensor.ID})
}

func (u *User) CreateSensorToUser(ctx context.Context, userID int64, sensor *domain.Sensor) error {
	if err := u.SensorRepo.SaveSensor(ctx, sensor); err != nil {
		return err
	}

	if _, err := u.UserRepo.GetUserByID(ctx, userID); err != nil {
		return err
	}

	return u.SensorOwnerRepo.SaveSensorOwner(ctx, domain.SensorOwner{UserID: userID, SensorID: sensor.ID})
}

func (u *User) GetUserSensors(ctx context.Context, userID int64) ([]domain.Sensor, error) {
	if _, err := u.UserRepo.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}

	sensorOwners, err := u.SensorOwnerRepo.GetSensorsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sensors := make([]domain.Sensor, len(sensorOwners))
	for i, so := range sensorOwners {
		s, err := u.SensorRepo.GetSensorByID(ctx, so.SensorID)
		if err != nil {
			return nil, err
		}
		sensors[i] = *s
	}
	return sensors, nil
}
