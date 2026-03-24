package clientcrypto

import (
	"bytes"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	key1 := DeriveKey([]byte("password"), []byte("salt"))
	key2 := DeriveKey([]byte("password"), []byte("salt"))
	key3 := DeriveKey([]byte("password"), []byte("different_salt"))
	key4 := DeriveKey([]byte("different_password"), []byte("salt"))

	if !bytes.Equal(key1, key2) {
		t.Fatal("same password and salt should produce same key")
	}

	if bytes.Equal(key1, key3) {
		t.Fatal("different salt should produce different key")
	}

	if bytes.Equal(key1, key4) {
		t.Fatal("different password should produce different key")
	}

	if len(key1) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key1))
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := DeriveKey([]byte("master_password"), []byte("static_salt"))
	data := []byte("sensitive data")

	encrypted, err := Encrypt(key, data)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if len(encrypted) <= len(data) {
		t.Fatal("encrypted data should be longer than plaintext (nonce + ciphertext + tag)")
	}

	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, data) {
		t.Fatalf("decrypted data mismatch: expected %q, got %q", data, decrypted)
	}
}

func TestEncryptDifferentNonce(t *testing.T) {
	key := DeriveKey([]byte("password"), []byte("salt"))
	data := []byte("same data")

	enc1, _ := Encrypt(key, data)
	enc2, _ := Encrypt(key, data)

	if bytes.Equal(enc1, enc2) {
		t.Fatal("encryption should use different nonce each time")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key1 := DeriveKey([]byte("password1"), []byte("salt"))
	key2 := DeriveKey([]byte("password2"), []byte("salt"))
	data := []byte("secret")

	encrypted, _ := Encrypt(key1, data)
	_, err := Decrypt(key2, encrypted)

	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptWithTamperedData(t *testing.T) {
	key := DeriveKey([]byte("password"), []byte("salt"))
	data := []byte("original")

	encrypted, _ := Encrypt(key, data)

	// Tamper with ciphertext
	encrypted[5] ^= 0xFF

	_, err := Decrypt(key, encrypted)
	if err == nil {
		t.Fatal("expected error when decrypting tampered data")
	}
}

func TestEmptyData(t *testing.T) {
	key := DeriveKey([]byte("password"), []byte("salt"))
	data := []byte("")

	encrypted, err := Encrypt(key, data)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if len(decrypted) != 0 {
		t.Fatalf("expected empty decrypted data, got %q", decrypted)
	}
}

func TestLargeData(t *testing.T) {
	key := DeriveKey([]byte("password"), []byte("salt"))
	data := make([]byte, 1024*1024) // 1MB
	for i := range data {
		data[i] = byte(i % 256)
	}

	encrypted, err := Encrypt(key, data)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, data) {
		t.Fatal("large data mismatch")
	}
}
