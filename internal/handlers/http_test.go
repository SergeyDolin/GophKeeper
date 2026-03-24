package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	h := NewHandler()

	body := []byte(`{"login":"test","password":"123"}`)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("expected 200")
	}
}

func TestRegisterDuplicateUser(t *testing.T) {
	h := NewHandler()

	body := []byte(`{"login":"dupuser","password":"123"}`)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("expected 200 for first registration")
	}

	// Try to register same user again
	req2 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w2 := httptest.NewRecorder()

	h.Router().ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Fatal("expected 400 for duplicate registration")
	}
}

func TestLogin(t *testing.T) {
	h := NewHandler()

	// First register a user
	regBody := []byte(`{"login":"logintest","password":"pass123"}`)
	regReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(regBody))
	regW := httptest.NewRecorder()
	h.Router().ServeHTTP(regW, regReq)

	// Now login
	loginBody := []byte(`{"login":"logintest","password":"pass123"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()

	h.Router().ServeHTTP(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", loginW.Code)
	}

	if len(loginW.Body.Bytes()) == 0 {
		t.Fatal("expected token in response")
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	h := NewHandler()

	// Register a user
	regBody := []byte(`{"login":"invuser","password":"pass123"}`)
	regReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(regBody))
	regW := httptest.NewRecorder()
	h.Router().ServeHTTP(regW, regReq)

	// Try login with wrong password
	loginBody := []byte(`{"login":"invuser","password":"wrongpass"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()

	h.Router().ServeHTTP(loginW, loginReq)

	if loginW.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", loginW.Code)
	}
}

func TestLoginNonExistentUser(t *testing.T) {
	h := NewHandler()

	loginBody := []byte(`{"login":"nouser","password":"pass123"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()

	h.Router().ServeHTTP(loginW, loginReq)

	if loginW.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", loginW.Code)
	}
}

func TestSecretsCreateAndGet(t *testing.T) {
	h := NewHandler()

	// Register and login to get a token
	regBody := []byte(`{"login":"secretsuser","password":"pass123"}`)
	regReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(regBody))
	regW := httptest.NewRecorder()
	h.Router().ServeHTTP(regW, regReq)

	loginBody := []byte(`{"login":"secretsuser","password":"pass123"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()
	h.Router().ServeHTTP(loginW, loginReq)

	token := loginW.Body.String()

	// Create a secret
	secretBody := []byte(`{"type":"password","data":"mysecretdata","meta":"mymeta"}`)
	secretReq := httptest.NewRequest(http.MethodPost, "/secrets", bytes.NewReader(secretBody))
	secretReq.Header.Set("Authorization", token)
	secretW := httptest.NewRecorder()

	h.Router().ServeHTTP(secretW, secretReq)

	if secretW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", secretW.Code)
	}

	// List secrets
	listReq := httptest.NewRequest(http.MethodGet, "/secrets", nil)
	listReq.Header.Set("Authorization", token)
	listW := httptest.NewRecorder()

	h.Router().ServeHTTP(listW, listReq)

	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listW.Code)
	}
}

func TestSecretsUnauthorized(t *testing.T) {
	h := NewHandler()

	// Try to access secrets without token
	secretReq := httptest.NewRequest(http.MethodGet, "/secrets", nil)
	secretW := httptest.NewRecorder()

	h.Router().ServeHTTP(secretW, secretReq)

	if secretW.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", secretW.Code)
	}
}

func TestSecretsDelete(t *testing.T) {
	h := NewHandler()

	// Register and login
	regBody := []byte(`{"login":"deluser","password":"pass123"}`)
	regReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(regBody))
	regW := httptest.NewRecorder()
	h.Router().ServeHTTP(regW, regReq)

	loginBody := []byte(`{"login":"deluser","password":"pass123"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()
	h.Router().ServeHTTP(loginW, loginReq)

	token := loginW.Body.String()

	// Create a secret
	secretBody := []byte(`{"type":"password","data":"todelete","meta":""}`)
	secretReq := httptest.NewRequest(http.MethodPost, "/secrets", bytes.NewReader(secretBody))
	secretReq.Header.Set("Authorization", token)
	secretW := httptest.NewRecorder()
	h.Router().ServeHTTP(secretW, secretReq)

	// Delete the secret (we need to get the ID from the create response)
	deleteReq := httptest.NewRequest(http.MethodDelete, "/secrets?id=1", nil)
	deleteReq.Header.Set("Authorization", token)
	deleteW := httptest.NewRecorder()

	h.Router().ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", deleteW.Code)
	}
}

func TestNewHandler(t *testing.T) {
	h := NewHandler()

	if h == nil {
		t.Fatal("expected non-nil handler")
	}

	if h.auth == nil {
		t.Fatal("expected auth service to be initialized")
	}

	if h.secrets == nil {
		t.Fatal("expected secret service to be initialized")
	}
}

func TestRouter(t *testing.T) {
	h := NewHandler()

	router := h.Router()

	if router == nil {
		t.Fatal("expected non-nil router")
	}
}
