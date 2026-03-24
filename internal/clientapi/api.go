// Package clientapi provides HTTP client for interacting with server API.
package clientapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const baseURL = "http://localhost:8080"

// Login authenticates user and returns JWT token.
func Login(login, password string) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"login":    login,
		"password": password,
	})

	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	token, _ := io.ReadAll(resp.Body)
	return string(token), nil
}

// AddSecret sends encrypted secret to server.
func AddSecret(token, typ, meta string, encrypted []byte) error {
	body, _ := json.Marshal(map[string]interface{}{
		"type": typ,
		"data": encrypted,
		"meta": meta,
	})

	req, _ := http.NewRequest("POST", baseURL+"/secrets", bytes.NewReader(body))
	req.Header.Set("Authorization", token)

	_, err := http.DefaultClient.Do(req)
	return err
}

// ListSecrets retrieves all user secrets from server.
func ListSecrets(token, password string) ([]map[string]interface{}, error) {
	body, _ := json.Marshal(map[string]string{
		"password": password,
	})
	req, _ := http.NewRequest("GET", baseURL+"/secrets", bytes.NewReader(body))
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}
