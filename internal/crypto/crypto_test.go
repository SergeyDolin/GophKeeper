package crypto

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	key := DeriveKey([]byte("pass"), []byte("salt"))

	data := []byte("hello")

	enc, _ := Encrypt(key, data)
	dec, _ := Decrypt(key, enc)

	if string(dec) != "hello" {
		t.Fatal("mismatch")
	}
}
