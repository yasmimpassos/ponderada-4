package handler

import (
	"encoding/json"
	"net/http"

	"racha-historico/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "nome, email e senha são obrigatórios")
		return
	}

	user, err := h.authService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusConflict, "email já cadastrado")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email e senha são obrigatórios")
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "credenciais inválidas")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
