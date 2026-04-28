package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrQuoteBelongsToAnotherUser = errors.New("unauthorized: quote belongs to another user")
	ErrQuoteExpired              = errors.New("this conversion quote has expired")
	ErrQuoteNotFound             = errors.New("quote not found")
	ErrRateFetch                 = errors.New("failed to fetch exchange rate")
)

type Quote struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SourceCurrency Currency
	TargetCurrency Currency
	ProviderRate   float64
	ExchangeRate   float64
	AmountIn       int64
	AmountOut      int64
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

type FXRateProvider interface {
	GetRate(ctx context.Context, source, target Currency) (float64, error)
}
