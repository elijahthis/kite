package interfaces

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/elijahthis/kite/internal/application"
	"github.com/elijahthis/kite/internal/domain"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	service *application.AuthService
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

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewAuthHandler(s application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: &s,
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Unable to parse request body", err)
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
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Unable to parse request body", err)
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

func writeError(w http.ResponseWriter, status int, errCode, msg string, err error) {
	log.Error().Err(err).Msg(err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errCode, Message: msg})
}
