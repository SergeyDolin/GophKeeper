package config

import (
	"os"
)

const fileName = ".gophkeeper_token"

// SaveToken - save JWT
func SaveToken(token string) error {
	return os.WriteFile(fileName, []byte(token), 0600)
}

// GetToken - read JWT
func GetToken() (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
