package service

import (
	"gophkeeper/internal/crypto"
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
func (s *SecretService) Create(user string, key []byte, typ, meta string, data []byte) (models.Secret, error) {
	encrypted, err := crypto.Encrypt(key, data)
	if err != nil {
		return models.Secret{}, err
	}

	secret := models.Secret{
		ID:        uuid.New().String(),
		UserLogin: user,
		Type:      typ,
		Data:      encrypted,
		Meta:      meta,
		UpdatedAt: time.Now(),
	}

	s.store.SaveSecret(secret)
	return secret, nil
}

// List - return secrets of user
func (s *SecretService) List(user string, key []byte) ([]models.Secret, error) {
	secrets := s.store.GetSecretByUser(user)

	for i := range secrets {
		decrypted, err := crypto.Decrypt(key, secrets[i].Data)
		if err != nil {
			return nil, err
		}
		secrets[i].Data = decrypted
	}

	return secrets, nil
}

// Delete - delete of secret
func (s *SecretService) Delete(id string) {
	s.store.DeleteSecret(id)
}
