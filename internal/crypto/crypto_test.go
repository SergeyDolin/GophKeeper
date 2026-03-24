package crypto

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	key := DeriveKey([]byte("password"), []byte("salt"))

	data := []byte("secret data")

	enc, err := Encrypt(key, data)
	if err != nil {
		t.Fatal(err)
	}

	dec, err := Decrypt(key, enc)
	if err != nil {
		t.Fatal(err)
	}

	if string(dec) != string(data) {
		t.Fatal("decrypt mismatch")
	}
}
