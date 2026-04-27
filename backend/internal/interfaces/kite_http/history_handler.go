package interfaces

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/elijahthis/kite/internal/application"
	"github.com/google/uuid"
)

type HistoryHandler struct {
	service *application.HistoryService
}

func NewHistoryHandler(s application.HistoryService) *HistoryHandler {
	return &HistoryHandler{service: &s}
}

func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusInternalServerError, "server_error", "Failed to identify user", nil)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	history, err := h.service.GetUserHistory(r.Context(), userID, page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fetch_failed", "Failed to retrieve transaction history", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"transactions": history,
		"page":         page,
		"limit":        limit,
	})
}
