package service

import (
	"main/internal/storage"
	"testing"
)

func TestRegisterAndLogin(t *testing.T) {
	store := storage.New()
	auth := NewAuthService(store)

	login := "user"
	password := "pass"

	if err := auth.Register(login, password); err != nil {
		t.Fatal(err)
	}

	token, err := auth.Login(login, password)
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("empty token")
	}
}

func TestInvalidPassword(t *testing.T) {
	store := storage.New()
	auth := NewAuthService(store)

	auth.Register("user", "pass")

	_, err := auth.Login("user", "wrong")
	if err == nil {
		t.Fatal("expected error")
	}
}
