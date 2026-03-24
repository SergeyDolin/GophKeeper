package config

import (
	"testing"
)

const testTokenFile = ".test_gophkeeper_token"

func TestSaveAndGetToken(t *testing.T) {

	token := "test.jwt.token"

	err := SaveToken(token)
	if err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	retrieved, err := GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if retrieved != token {
		t.Fatalf("token mismatch: expected %s, got %s", token, retrieved)
	}
}

func TestSaveTokenOverwrite(t *testing.T) {

	token1 := "first.token"
	token2 := "second.token"

	if err := SaveToken(token1); err != nil {
		t.Fatalf("first SaveToken failed: %v", err)
	}

	if err := SaveToken(token2); err != nil {
		t.Fatalf("second SaveToken failed: %v", err)
	}

	retrieved, err := GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if retrieved != token2 {
		t.Fatalf("expected overwritten token %s, got %s", token2, retrieved)
	}
}

func TestSaveEmptyToken(t *testing.T) {

	err := SaveToken("")
	if err != nil {
		t.Fatalf("SaveToken with empty string failed: %v", err)
	}

	retrieved, err := GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if retrieved != "" {
		t.Fatalf("expected empty token, got %s", retrieved)
	}
}

func TestTokenWithSpecialCharacters(t *testing.T) {

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	err := SaveToken(token)
	if err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	retrieved, err := GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if retrieved != token {
		t.Fatal("JWT token mismatch")
	}
}
