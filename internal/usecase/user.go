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

func (u *User) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := u.UserRepo.SaveUser(ctx, user); err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) AttachSensorToUser(ctx context.Context, userID, sensorID int64) error {
	user, err := u.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	sensor, err := u.SensorRepo.GetSensorByID(ctx, sensorID)
	if err != nil {
		return err
	}

	return u.SensorOwnerRepo.SaveSensorOwner(ctx, domain.SensorOwner{UserID: user.ID, SensorID: sensor.ID})
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