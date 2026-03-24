package storage

import (
	"gophkeeper/internal/models"
	"sync"
	"testing"
)

func TestNewStorage(t *testing.T) {
	store := New()
	if store == nil {
		t.Fatal("New() returned nil")
	}
	if store.users == nil {
		t.Fatal("users map not initialized")
	}
}

func TestSaveUser(t *testing.T) {
	store := New()
	user := models.User{
		Login:        "testuser",
		PasswordHash: []byte("hash"),
		Salt:         []byte("salt"),
	}

	err := store.SaveUser(user)
	if err != nil {
		t.Fatalf("SaveUser failed: %v", err)
	}

	retrieved, err := store.GetUser("testuser")
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	if retrieved.Login != user.Login {
		t.Fatalf("login mismatch: expected %s, got %s", user.Login, retrieved.Login)
	}
}

func TestSaveUserDuplicate(t *testing.T) {
	store := New()
	user := models.User{
		Login:        "testuser",
		PasswordHash: []byte("hash"),
		Salt:         []byte("salt"),
	}

	if err := store.SaveUser(user); err != nil {
		t.Fatalf("first SaveUser failed: %v", err)
	}

	err := store.SaveUser(user)
	if err == nil {
		t.Fatal("expected error on duplicate user")
	}
}

func TestGetUserNotFound(t *testing.T) {
	store := New()
	_, err := store.GetUser("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestSaveSecret(t *testing.T) {
	store := New()
	secret := models.Secret{
		ID:        "secret-1",
		UserLogin: "testuser",
		Type:      "text",
		Data:      []byte("encrypted data"),
		Meta:      "metadata",
	}

	store.SaveSecret(secret)

	secrets := store.GetSecretByUser("testuser")
	if len(secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(secrets))
	}

	if secrets[0].ID != secret.ID {
		t.Fatalf("secret ID mismatch")
	}
}

func TestGetSecretByUserMultiple(t *testing.T) {
	store := New()

	secrets := []models.Secret{
		{ID: "1", UserLogin: "user1", Type: "text", Data: []byte("data1")},
		{ID: "2", UserLogin: "user1", Type: "login", Data: []byte("data2")},
		{ID: "3", UserLogin: "user2", Type: "card", Data: []byte("data3")},
	}

	for _, s := range secrets {
		store.SaveSecret(s)
	}

	user1Secrets := store.GetSecretByUser("user1")
	if len(user1Secrets) != 2 {
		t.Fatalf("expected 2 secrets for user1, got %d", len(user1Secrets))
	}

	user2Secrets := store.GetSecretByUser("user2")
	if len(user2Secrets) != 1 {
		t.Fatalf("expected 1 secret for user2, got %d", len(user2Secrets))
	}

	user3Secrets := store.GetSecretByUser("user3")
	if len(user3Secrets) != 0 {
		t.Fatalf("expected 0 secrets for user3, got %d", len(user3Secrets))
	}
}

func TestDeleteSecret(t *testing.T) {
	store := New()
	secret := models.Secret{
		ID:        "to-delete",
		UserLogin: "testuser",
		Type:      "text",
		Data:      []byte("data"),
	}

	store.SaveSecret(secret)

	secrets := store.GetSecretByUser("testuser")
	if len(secrets) != 1 {
		t.Fatalf("expected 1 secret before delete")
	}

	store.DeleteSecret("to-delete")

	secrets = store.GetSecretByUser("testuser")
	if len(secrets) != 0 {
		t.Fatalf("expected 0 secrets after delete, got %d", len(secrets))
	}
}

func TestDeleteNonExistentSecret(t *testing.T) {
	store := New()
	store.DeleteSecret("nonexistent") // Should not panic
}

func TestConcurrentAccess(t *testing.T) {
	store := New()
	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(id int) {
			defer wg.Done()
			user := models.User{
				Login:        "user",
				PasswordHash: []byte("hash"),
				Salt:         []byte("salt"),
			}
			store.SaveUser(user)
		}(i)
	}

	wg.Wait()
}

func TestConcurrentSecretAccess(t *testing.T) {
	store := New()
	var wg sync.WaitGroup

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(id int) {
			defer wg.Done()
			secret := models.Secret{
				ID:        "secret",
				UserLogin: "user",
				Type:      "text",
				Data:      []byte("data"),
			}
			store.SaveSecret(secret)
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(id int) {
			defer wg.Done()
			store.GetSecretByUser("user")
		}(i)
	}

	wg.Wait()
}
