package service

import (
	"gophkeeper/internal/storage"
	"testing"
)

func TestSecretCRUD(t *testing.T) {
	store := storage.New()
	service := NewSecretService(store)

	user := "test"

	// create
	sec := service.Create(user, "text", "meta", []byte("hello"))

	if sec.ID == "" {
		t.Fatal("empty id")
	}

	// list
	list := service.List(user)
	if len(list) != 1 {
		t.Fatal("expected 1 secret")
	}

	// delete
	service.Delete(sec.ID)

	list = service.List(user)
	if len(list) != 0 {
		t.Fatal("delete failed")
	}
}
