package service

import (
	"gophkeeper/internal/models"
	"gophkeeper/internal/storage"
	"time"

	"github.com/google/uuid"
)

// SecretService manager user secrets.
type SecretService struct {
	store *storage.MemoryStorage
}

// NewSecretService creates a new SecretService instance.
func NewSecretService(store *storage.MemoryStorage) *SecretService {
	return &SecretService{store: store}
}

// Create stores a new encrypted secret for the user.
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

// List returns all secrets for a given user.
func (s *SecretService) List(user string) []models.Secret {
	return s.store.GetSecretByUser(user)
}

// Delete removes a secret by ID.
func (s *SecretService) Delete(id string) {
	s.store.DeleteSecret(id)
}
