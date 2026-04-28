package interfaces

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.
			Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("ip", r.RemoteAddr).
			Msgf("Incoming Request: ")
		next.ServeHTTP(w, r)
	})
}

// will add logging middleware later

func RequireAuth(secretKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("kite_session")
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "Missing session cookie", err)
				return
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "unauthorized", "Invalid or expired session", err)
				return
			}

			subject, err := token.Claims.GetSubject()
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "Missing Subject", err)
				return
			}
			userID, err := uuid.Parse(subject)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "Invalid user ID in token", err)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
