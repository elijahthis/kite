package interfaces

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/elijahthis/kite/internal/application"
	"github.com/elijahthis/kite/internal/domain"
)

type AuthHandler struct {
	service application.AuthService
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func NewAuthHandler(s application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	req, ok := parseJSON[RegisterRequest](w, r)
	if !ok {
		return
	}

	user, err := ah.service.Register(r.Context(), req.Email, req.Password, req.FirstName)
	if err != nil {
		if err == domain.ErrUserAlreadyExists {
			writeError(w, http.StatusBadRequest, "signup_failed", err.Error(), err)
			return
		}
		writeError(w, http.StatusInternalServerError, "server_error", "An internal error occurred", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"ID":     user.ID.String(),
	})
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, ok := parseJSON[LoginRequest](w, r)
	if !ok {
		return
	}

	token, err := ah.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			writeError(w, http.StatusBadRequest, "auth_failed", err.Error(), err)
			return
		}
		writeError(w, http.StatusInternalServerError, "server_error", "An internal error occurred", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "kite_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
