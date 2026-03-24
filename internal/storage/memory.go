package storage

import (
	"errors"
	"gophkeeper/internal/models"
	"sync"
)

type MemoryStorage struct {
	mu      sync.RWMutex
	users   map[string]models.User
	secrets map[string]models.Secret
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

// SaveSecret - save secret
func (s *MemoryStorage) SaveSecret(secret models.Secret) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.secrets == nil {
		s.secrets = make(map[string]models.Secret)
	}

	s.secrets[secret.ID] = secret
}

// GetSecretByUser - return all secrets of users
func (s *MemoryStorage) GetSecretByUser(login string) []models.Secret {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Secret

	for _, sec := range s.secrets {
		if sec.UserLogin == login {
			result = append(result, sec)
		}
	}

	return result
}

// DeleteSecret - delete secret
func (s *MemoryStorage) DeleteSecret(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.secrets, id)
}
