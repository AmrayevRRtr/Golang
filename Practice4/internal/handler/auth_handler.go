package handler

import (
	"Practice4/internal/auth"
	"Practice4/internal/usecase"
	"encoding/json"
	"net/http"
)

const (
	AdminEmail    = "admin@mail.com"
	AdminPassword = "password123"
	AdminID       = 22
)

type AuthHandler struct {
	usecase *usecase.UserUsecase
}

func NewAuthHandler(usecase *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email != AdminEmail || req.Password != AdminPassword {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(AdminID, AdminEmail, "admin")
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
