package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/elijahthis/kite/internal/application"
	"github.com/elijahthis/kite/internal/domain"
)

type PayoutHandler struct {
	service application.PayoutService
}

func NewPayoutHandler(s application.PayoutService) *PayoutHandler {
	return &PayoutHandler{service: s}
}

type PayoutRequest struct {
	SourceCurrency string `json:"source_currency"`
	Amount         int64  `json:"amount"`
	AccountNumber  string `json:"account_number"`
	BankCode       string `json:"bank_code"`
	AccountName    string `json:"account_name"`
}

func (ph *PayoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(w, r)
	if !ok {
		return
	}

	req, ok := parseJSON[PayoutRequest](w, r)
	if !ok {
		return
	}

	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_amount", "Payout amount must be strictly positive", nil)
		return
	}
	if req.AccountNumber == "" || req.BankCode == "" || req.AccountName == "" {
		writeError(w, http.StatusBadRequest, "missing_bank_details", "Recipient bank details are required", nil)
		return
	}

	parsedCurrency, ok := domain.GetCurrency(req.SourceCurrency)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_currency", "The provided source currency is invalid", nil)
		return
	}

	txn, err := ph.service.ExecutePayout(r.Context(), userID, parsedCurrency, req.Amount, req.AccountNumber, req.BankCode)
	if err != nil {
		if err == application.ErrInsufficientFunds {
			writeError(w, http.StatusBadRequest, "insufficient_funds", err.Error(), nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "payout_failed", "Failed to initiate payout", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "pending",
		"message":        "Payout initiated and is processing",
		"transaction_id": txn.ID.String(),
	})
}
