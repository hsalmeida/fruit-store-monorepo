package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hsalmeida/fruit-store-monorepo/api/internal/auth"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/repository"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginHandler struct {
	svc service.UserService
}

func NewLoginHandler(db *pgxpool.Pool) *LoginHandler {
	svc := service.NewUserService(repository.NewUserRepository(db))
	return &LoginHandler{svc: svc}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	user, err := h.svc.Authenticate(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
