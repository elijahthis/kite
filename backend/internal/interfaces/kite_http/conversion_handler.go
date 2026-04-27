package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/elijahthis/kite/internal/application"
	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
)

type ConversionHandler struct {
	service *application.ConversionService
}

func NewConversionHandler(s application.ConversionService) *ConversionHandler {
	return &ConversionHandler{service: &s}
}

type QuoteRequest struct {
	SourceCurrency string `json:"source_currency"`
	TargetCurrency string `json:"target_currency"`
	AmountIn       int64  `json:"amount_in"`
}

type ExecuteRequest struct {
	QuoteID string `json:"quote_id"`
}

func (h *ConversionHandler) GenerateQuote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusInternalServerError, "server_error", "Failed to identify user", nil)
		return
	}

	var req QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Unable to parse request body", err)
		return
	}

	src, okSrc := domain.GetCurrency(req.SourceCurrency)
	tgt, okTgt := domain.GetCurrency(req.TargetCurrency)
	if !okSrc || !okTgt || src == tgt {
		writeError(w, http.StatusBadRequest, "invalid_currency", "Invalid or identical currencies", nil)
		return
	}

	quote, err := h.service.GenerateQuote(r.Context(), userID, src, tgt, req.AmountIn)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "quote_failed", "Failed to generate quote", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"quote_id":      quote.ID,
		"exchange_rate": quote.ExchangeRate,
		"amount_in":     quote.AmountIn,
		"amount_out":    quote.AmountOut,
		"expires_at":    quote.ExpiresAt,
	})
}

func (h *ConversionHandler) ExecuteQuote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusInternalServerError, "server_error", "Failed to identify user", nil)
		return
	}

	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Unable to parse request body", err)
		return
	}

	quoteID, err := uuid.Parse(req.QuoteID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_quote_id", "Malformed quote ID", err)
		return
	}

	err = h.service.ExecuteQuote(r.Context(), userID, quoteID)
	if err != nil {
		if err == domain.ErrQuoteExpired {
			writeError(w, http.StatusBadRequest, "quote_expired", err.Error(), nil)
			return
		}
		if err == application.ErrInsufficientFunds {
			writeError(w, http.StatusBadRequest, "insufficient_funds", err.Error(), nil)
			return
		}
		if err == domain.ErrDuplicateReference {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Quote already executed"})
			return
		}
		writeError(w, http.StatusInternalServerError, "execution_failed", "Failed to execute conversion", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
