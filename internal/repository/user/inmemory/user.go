package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"math/rand"
	"sync"
)

var (
	ErrUserNil      = errors.New("user is nil")
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	mu          sync.RWMutex
	UserStorage map[int64]domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{UserStorage: make(map[int64]domain.User)}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if user == nil {
		return ErrUserNil
	}

	user.ID = rand.Int63()
	r.mu.Lock()
	r.UserStorage[user.ID] = *user
	r.mu.Unlock()
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	r.mu.RLock()
	if usr, ok := r.UserStorage[id]; ok {
		return &usr, nil
	}
	r.mu.Unlock()
	return nil, ErrUserNotFound
}