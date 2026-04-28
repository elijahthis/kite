package interfaces

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type contextKey string
type cookieKey string

const UserIDKey contextKey = "user_id"
const kiteSessionKey cookieKey = "kite_session"

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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

func RequireAuth(secretKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(string(kiteSessionKey))
			if err != nil {
				writeError(r.Context(), w, http.StatusUnauthorized, "unauthorized", "Missing session cookie", err)
				return
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				writeError(r.Context(), w, http.StatusUnauthorized, "unauthorized", "Invalid or expired session", err)
				return
			}

			subject, err := token.Claims.GetSubject()
			if err != nil {
				writeError(r.Context(), w, http.StatusUnauthorized, "unauthorized", "Missing Subject", err)
				return
			}
			userID, err := uuid.Parse(subject)
			if err != nil {
				writeError(r.Context(), w, http.StatusUnauthorized, "unauthorized", "Invalid user ID in token", err)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func RequestLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", reqID)

		reqLogger := log.With().Str("request_id", reqID).Logger()

		ctx := reqLogger.WithContext(r.Context())
		r = r.WithContext(ctx)

		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		reqLogger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Msg("Request started")

		next.ServeHTTP(rec, r)

		reqLogger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rec.statusCode).
			Dur("latency", time.Since(start)).
			Msg("Request completed")
	}
}

func ApplyMiddleware(h http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
