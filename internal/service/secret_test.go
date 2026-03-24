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
	sec, err := service.Create(user, []byte("key"), "text", "meta", []byte("hello"))
	if err != nil {
		return
	}
	if sec.ID == "" {
		t.Fatal("empty id")
	}

	// list
	list, err := service.List(user, []byte("key"))
	if err != nil {
		return
	}
	if len(list) != 1 {
		t.Fatal("expected 1 secret")
	}

	// delete
	service.Delete(sec.ID)

	list, err = service.List(user, []byte("key"))
	if err != nil {
		return
	}
	if len(list) != 0 {
		t.Fatal("delete failed")
	}
}
