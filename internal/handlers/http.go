package handlers

import (
	"encoding/json"
	"main/internal/service"
	"main/internal/storage"
	"net/http"
)

type Handler struct {
	auth *service.AuthService
}

func NewHandler() *Handler {
	store := storage.New()
	return &Handler{
		auth: service.NewAuthService(store),
	}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", h.register)
	mux.HandleFunc("/login", h.login)

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
