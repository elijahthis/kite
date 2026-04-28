package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/elijahthis/kite/internal/application"
)

type WalletHandler struct {
	service application.WalletService
}

func NewWalletHandler(s application.WalletService) *WalletHandler {
	return &WalletHandler{service: s}
}

func (h *WalletHandler) GetBalances(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(w, r)
	if !ok {
		return
	}

	balances, err := h.service.GetBalances(r.Context(), userID)
	if err != nil {
		writeError(r.Context(), w, http.StatusInternalServerError, "fetch_failed", "Failed to retrieve balances", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"balances": balances,
	})
}
