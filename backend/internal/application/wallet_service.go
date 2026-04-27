package application

import (
	"context"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
)

type WalletService struct {
	ledgerRepo domain.LedgerRepository
}

func NewWalletService(lr domain.LedgerRepository) *WalletService {
	return &WalletService{ledgerRepo: lr}
}

func (s *WalletService) GetBalances(ctx context.Context, userID uuid.UUID) (map[string]int64, error) {
	activeBalances, err := s.ledgerRepo.GetAllAccountBalances(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := map[string]int64{
		domain.USD.String(): 0,
		domain.EUR.String(): 0,
		domain.GBP.String(): 0,
		domain.NGN.String(): 0,
		domain.KES.String(): 0,
	}

	for currency, balance := range activeBalances {
		response[currency.String()] = balance
	}

	return response, nil
}
