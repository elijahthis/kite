package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func getUserIDFromContext(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		writeError(r.Context(), w, http.StatusInternalServerError, "server_error", "Failed to identify user", nil)
		return uuid.UUID{}, false
	}
	return userID, ok
}

func parseJSON[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(r.Context(), w, http.StatusBadRequest, "invalid_request", "Unable to parse request body", err)
		return req, false
	}
	return req, true
}
