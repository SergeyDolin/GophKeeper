package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := DeriveKey([]byte("pass"), []byte("salt"))

	data := []byte("hello")

	enc, _ := Encrypt(key, data)
	dec, _ := Decrypt(key, enc)

	if string(dec) != "hello" {
		t.Fatal("mismatch")
	}
}

func TestDeriveKey(t *testing.T) {
	password := []byte("mypassword")
	salt := []byte("randomsalt")

	key1 := DeriveKey(password, salt)
	key2 := DeriveKey(password, salt)

	if !bytes.Equal(key1, key2) {
		t.Fatal("derived keys should be equal for same password and salt")
	}

	if len(key1) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key1))
	}

	// Different salt should produce different key
	key3 := DeriveKey(password, []byte("differentsalt"))
	if bytes.Equal(key1, key3) {
		t.Fatal("different salts should produce different keys")
	}

	// Different password should produce different key
	key4 := DeriveKey([]byte("differentpassword"), salt)
	if bytes.Equal(key1, key4) {
		t.Fatal("different passwords should produce different keys")
	}
}

func TestEncryptDifferentNonces(t *testing.T) {
	key := DeriveKey([]byte("pass"), []byte("salt"))
	data := []byte("test data")

	// Encrypt the same data twice
	enc1, err := Encrypt(key, data)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	enc2, err := Encrypt(key, data)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Ciphertexts should be different due to random nonce
	if bytes.Equal(enc1, enc2) {
		t.Fatal("encrypting the same data twice should produce different ciphertexts")
	}

	// Both should decrypt to the same plaintext
	dec1, err := Decrypt(key, enc1)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	dec2, err := Decrypt(key, enc2)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(dec1, dec2) {
		t.Fatal("both decryptions should produce the same plaintext")
	}
}

func TestEncryptEmptyData(t *testing.T) {
	key := DeriveKey([]byte("pass"), []byte("salt"))
	data := []byte("")

	enc, err := Encrypt(key, data)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	dec, err := Decrypt(key, enc)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if len(dec) != 0 {
		t.Fatalf("expected empty decrypted data, got %d bytes", len(dec))
	}
}

func TestEncryptLargeData(t *testing.T) {
	key := DeriveKey([]byte("pass"), []byte("salt"))

	// Create large data (1MB)
	data := make([]byte, 1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	enc, err := Encrypt(key, data)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	dec, err := Decrypt(key, enc)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(data, dec) {
		t.Fatal("large data mismatch after encryption/decryption")
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	key := DeriveKey([]byte("pass"), []byte("salt"))

	// Corrupted ciphertext
	validData := []byte("test")
	validEnc, _ := Encrypt(key, validData)

	// Corrupt some bytes (skip the nonce which is at the beginning)
	nonceSize := 12 // AES-GCM standard nonce size
	corrupted := make([]byte, len(validEnc))
	copy(corrupted, validEnc)
	if len(corrupted) > nonceSize {
		corrupted[nonceSize] ^= 0xFF
	}

	_, err := Decrypt(key, corrupted)
	if err == nil {
		t.Fatal("expected error for corrupted ciphertext")
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := DeriveKey([]byte("pass1"), []byte("salt"))
	key2 := DeriveKey([]byte("pass2"), []byte("salt"))

	data := []byte("secret message")

	enc, _ := Encrypt(key1, data)

	// Try to decrypt with wrong key
	_, err := Decrypt(key2, enc)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestEncryptNilPlaintext(t *testing.T) {
	key := DeriveKey([]byte("pass"), []byte("salt"))

	enc, err := Encrypt(key, nil)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	dec, err := Decrypt(key, enc)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if dec != nil && len(dec) != 0 {
		t.Fatalf("expected nil or empty decrypted data, got %v", dec)
	}
}
