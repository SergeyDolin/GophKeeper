// Package handlers provides HTTP handlers for the API.
package handlers

import (
	"encoding/json"
	"gophkeeper/internal/service"
	"gophkeeper/internal/storage"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Handler handlers HTTP requests for authentication and secrets.
type Handler struct {
	auth    *service.AuthService
	secrets *service.SecretService
}

// NewHandler initializes all dependencies and returns a Handler.
func NewHandler() *Handler {
	store := storage.New()
	return &Handler{
		auth:    service.NewAuthService(store),
		secrets: service.NewSecretService(store),
	}
}

// Router returns configured HTTP routes.
func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", h.register)
	mux.HandleFunc("/login", h.login)
	mux.HandleFunc("/secrets", h.secretsHandler)

	return mux
}

type request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// register handles user registration.
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req request
	json.NewDecoder(r.Body).Decode(&req)

	err := h.auth.Register(req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// login handles user authentication and returns JWT.
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req request
	json.NewDecoder(r.Body).Decode(&req)

	token, err := h.auth.Login(req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Write([]byte(token))
}

type secretRequest struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
	Meta string `json:"meta"`
}

// secretsHandler handles CRUD operations for secrets.
func (h *Handler) secretsHandler(w http.ResponseWriter, r *http.Request) {
	user := h.getUserFromRequest(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		var req secretRequest
		json.NewDecoder(r.Body).Decode(&req)

		sec := h.secrets.Create(user, req.Type, req.Meta, req.Data)

		json.NewEncoder(w).Encode(sec)

	case http.MethodGet:
		list := h.secrets.List(user)
		json.NewEncoder(w).Encode(list)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		h.secrets.Delete(id)
		w.WriteHeader(http.StatusOK)
	}
}

// getUserFromRequest extracts user login from JWT token.
func (h *Handler) getUserFromRequest(r *http.Request) string {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return ""
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || token == nil {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	login, ok := claims["login"].(string)
	if !ok {
		return ""
	}

	return login
}
