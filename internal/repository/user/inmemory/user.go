package inmemory

import (
	"context"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
)

type UserRepository struct {
	mu          sync.RWMutex
	UserStorage []domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{UserStorage: make([]domain.User, 0)}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if user == nil {
		return usecase.ErrUserNotFound
	}

	r.mu.Lock()
	r.UserStorage = append(r.UserStorage, *user)
	r.mu.Unlock()
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.mu.RLock()
	for _, usr := range r.UserStorage {
		if usr.ID == id {
			return &usr, nil
		}
	}
	r.mu.RUnlock()
	return nil, usecase.ErrUserNotFound
}
