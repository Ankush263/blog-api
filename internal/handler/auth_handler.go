package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Ankush263/blog-api/internal/auth"
	"github.com/Ankush263/blog-api/internal/model"
	"github.com/Ankush263/blog-api/internal/repository"
)

type AuthHandler struct {
	users *repository.UserRepository
}

func NewAuthHandler(users *repository.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
}

type authRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "password error", http.StatusInternalServerError)
		return
	}

	user := model.User{
		Email: req.Email,
		Password: hash,
	}

	if err := h.users.Create(r.Context(), &user); err != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(int64(user.ID))
	if err != nil {
		http.Error(w, "Token Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}
