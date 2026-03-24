package handlers

import (
	"encoding/json"
	"gophkeeper/internal/service"
	"gophkeeper/internal/storage"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	auth    *service.AuthService
	secrets *service.SecretService
}

func NewHandler() *Handler {
	store := storage.New()
	return &Handler{
		auth:    service.NewAuthService(store),
		secrets: service.NewSecretService(store),
	}
}

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

func (h *Handler) secretsHandler(w http.ResponseWriter, r *http.Request) {
	user := h.getUserFromRequest(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		Type     string `json:"type"`
		Data     string `json:"data"`
		Meta     string `json:"meta"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	key, err := h.auth.GetKey(user, req.Password)
	if err != nil {
		http.Error(w, "invalid password", 401)
		return
	}

	switch r.Method {
	case http.MethodPost:
		sec, err := h.secrets.Create(user, key, req.Type, req.Meta, []byte(req.Data))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(sec)

	case http.MethodGet:
		list, err := h.secrets.List(user, key)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(list)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		h.secrets.Delete(id)
		w.WriteHeader(http.StatusOK)
	}
}

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
