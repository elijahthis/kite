package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/elijahthis/kite/internal/application"
	"github.com/elijahthis/kite/internal/domain"
)

type DepositHandler struct {
	service *application.DepositService
}

type DepositRequest struct {
	Currency  string `json:"currency"`
	Amount    int64  `json:"amount"`
	Reference string `json:"reference"`
}

func NewDepositHandler(s application.DepositService) *DepositHandler {
	return &DepositHandler{
		service: &s,
	}
}

func (dh *DepositHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(w, r)
	if !ok {
		return
	}

	req, ok := parseJSON[DepositRequest](w, r)
	if !ok {
		return
	}

	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_amount", "Deposit amount must be greater than zero", nil)
		return
	}

	if req.Reference == "" {
		writeError(w, http.StatusBadRequest, "missing_reference", "Idempotency reference is required", nil)
		return
	}
	parsedCurrency, ok := domain.GetCurrency(req.Currency)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_currency", "Currency parsed is invalid", nil)
		return
	}

	if err := dh.service.ExecuteDeposit(r.Context(), userID, parsedCurrency, req.Amount, req.Reference); err != nil {
		if err == domain.ErrDuplicateReference {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "success",
				"message": "Deposit already processed",
			})
			return
		}

		writeError(w, http.StatusInternalServerError, "deposit_failed", "Failed to process deposit", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Deposit successful",
	})
}
