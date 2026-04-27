package application

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
)

var ErrInsufficientFunds = errors.New("insufficient funds for conversion")

type ConversionService struct {
	fxProvider  domain.FXRateProvider
	quoteRepo   domain.QuoteRepository
	atomicUnit  domain.AtomicUnit
	userRepo    domain.UserRepository
	accountRepo domain.AccountRepository
	ledgerRepo  domain.LedgerRepository
	sysEmail    string
}

func NewConversionService(
	fx domain.FXRateProvider,
	qr domain.QuoteRepository,
	at domain.AtomicUnit,
	ur domain.UserRepository,
	ar domain.AccountRepository,
	lr domain.LedgerRepository,
	sysEmail string,
) *ConversionService {
	return &ConversionService{
		fxProvider:  fx,
		quoteRepo:   qr,
		atomicUnit:  at,
		userRepo:    ur,
		accountRepo: ar,
		ledgerRepo:  lr,
		sysEmail:    sysEmail,
	}
}

func (s *ConversionService) GenerateQuote(ctx context.Context, userID uuid.UUID, source, target domain.Currency, amountIn int64) (*domain.Quote, error) {

	providerRate, err := s.fxProvider.GetRate(ctx, source, target)
	if err != nil {
		return nil, err
	}

	spreadPercent := 0.01
	exchangeRate := providerRate * (1.0 - spreadPercent)

	amountOutFloat := float64(amountIn) * exchangeRate
	amountOut := int64(math.Round(amountOutFloat))

	now := time.Now().UTC()
	quote := &domain.Quote{
		ID:             uuid.New(),
		UserID:         userID,
		SourceCurrency: source,
		TargetCurrency: target,
		ProviderRate:   providerRate,
		ExchangeRate:   exchangeRate,
		AmountIn:       amountIn,
		AmountOut:      amountOut,
		ExpiresAt:      now.Add(60 * time.Second),
		CreatedAt:      now,
	}

	if err := s.quoteRepo.Save(ctx, quote); err != nil {
		return nil, err
	}

	return quote, nil
}

func (s *ConversionService) ExecuteQuote(ctx context.Context, userID uuid.UUID, quoteID uuid.UUID) error {
	return s.atomicUnit.Do(ctx, func(ctxWithTx context.Context) error {
		quote, err := s.quoteRepo.Get(ctxWithTx, quoteID)
		if err != nil {
			return err
		}

		if quote.UserID != userID {
			return errors.New("unauthorized: quote belongs to another user")
		}

		if time.Now().UTC().After(quote.ExpiresAt) {
			return domain.ErrQuoteExpired
		}

		sysUser, err := s.userRepo.FindByEmail(ctxWithTx, s.sysEmail)
		if err != nil {
			return fmt.Errorf("system user not found: %w", err)
		}

		userSourceAcct, err := getOrCreateAccountUtil(s.accountRepo, ctxWithTx, userID, quote.SourceCurrency)
		if err != nil {
			return err
		}

		userTargetAcct, err := getOrCreateAccountUtil(s.accountRepo, ctxWithTx, userID, quote.TargetCurrency)
		if err != nil {
			return err
		}

		sysSourceAcct, err := getOrCreateAccountUtil(s.accountRepo, ctxWithTx, sysUser.ID, quote.SourceCurrency)
		if err != nil {
			return err
		}

		sysTargetAcct, err := getOrCreateAccountUtil(s.accountRepo, ctxWithTx, sysUser.ID, quote.TargetCurrency)
		if err != nil {
			return err
		}

		currentBalance, err := s.ledgerRepo.GetAccountBalance(ctxWithTx, userSourceAcct)
		if err != nil {
			return err
		}
		if currentBalance < quote.AmountIn {
			return ErrInsufficientFunds
		}

		builder := domain.NewTransactionBuilder(domain.CONVERSION, domain.SUCCESS, quote.ID.String())

		// User to System
		builder.AddEntry(userSourceAcct.ID, quote.AmountIn, domain.DEBIT, quote.SourceCurrency)
		builder.AddEntry(sysSourceAcct.ID, quote.AmountIn, domain.CREDIT, quote.SourceCurrency)

		// System to User
		builder.AddEntry(sysTargetAcct.ID, quote.AmountOut, domain.DEBIT, quote.TargetCurrency)
		builder.AddEntry(userTargetAcct.ID, quote.AmountOut, domain.CREDIT, quote.TargetCurrency)

		txn, err := builder.Build()
		if err != nil {
			return fmt.Errorf("failed to build conversion txn: %w", err)
		}

		if err := s.ledgerRepo.AppendTransaction(ctxWithTx, txn); err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				return domain.ErrDuplicateReference
			}
			return err
		}

		return nil
	})
}
