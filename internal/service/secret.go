package service

import (
	"gophkeeper/internal/models"
	"gophkeeper/internal/storage"
	"time"

	"github.com/google/uuid"
)

type SecretService struct {
	store *storage.MemoryStorage
}

func NewSecretService(store *storage.MemoryStorage) *SecretService {
	return &SecretService{store: store}
}

// Create - create a secret
func (s *SecretService) Create(user, typ, meta string, data []byte) models.Secret {
	secret := models.Secret{
		ID:        uuid.New().String(),
		UserLogin: user,
		Type:      typ,
		Data:      data,
		Meta:      meta,
		UpdatedAt: time.Now(),
	}

	s.store.SaveSecret(secret)
	return secret
}

// List - return secrets of user
func (s *SecretService) List(user string) []models.Secret {
	return s.store.GetSecretByUser(user)
}

// Delete - delete of secret
func (s *SecretService) Delete(id string) {
	s.store.DeleteSecret(id)
}
