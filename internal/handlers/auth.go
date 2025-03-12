package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shortly/internal/models"
	"github.com/shortly/internal/services"
	"github.com/shortly/internal/utils"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		writeError(w, "username must be 3-50 chars", http.StatusBadRequest)
		return
	}
	if !utils.IsValidEmail(req.Email) {
		writeError(w, "invalid email", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		writeError(w, "password must be at least 6 chars", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Register(r.Context(), req)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, resp, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	writeJSON(w, resp, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
