package application

import (
	"context"
	"fmt"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
)

type Services struct {
	Auth       *AuthService
	Deposit    *DepositService
	Conversion *ConversionService
	Wallet     *WalletService
	Payout     *PayoutService
}

func getOrCreateAccountUtil(accountRepo domain.AccountRepository, ctx context.Context, ownerID uuid.UUID, currency domain.Currency) (*domain.Account, error) {
	account, err := accountRepo.GetByUserIDAndCurrency(ctx, ownerID, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to verify account existence: %w", err)
	}
	if account != nil {
		return account, nil
	}

	newAcct := &domain.Account{
		UserID:   ownerID,
		Currency: currency,
	}

	if err := accountRepo.Create(ctx, newAcct); err != nil {
		return nil, fmt.Errorf("failed to create new account: %w", err)
	}

	return newAcct, nil
}
