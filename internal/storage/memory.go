package storage

import (
	"errors"
	"main/internal/models"
	"sync"
)

type MemoryStorage struct {
	mu    sync.RWMutex
	users map[string]models.User
}

func New() *MemoryStorage {
	return &MemoryStorage{
		users: make(map[string]models.User),
	}
}

// SaveUser - save user information
func (s *MemoryStorage) SaveUser(u models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[u.Login]; exists {
		return errors.New("user already exists")
	}

	s.users[u.Login] = u
	return nil
}

// GetUser - get user information
func (s *MemoryStorage) GetUser(login string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[login]
	if !ok {
		return models.User{}, errors.New("user not found")
	}

	return u, nil
}
