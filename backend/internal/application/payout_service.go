package application

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PayoutService struct {
	atomicUnit  domain.AtomicUnit
	userRepo    domain.UserRepository
	accountRepo domain.AccountRepository
	ledgerRepo  domain.LedgerRepository
	sysEmail    string
}

func NewPayoutService(atomicUnit domain.AtomicUnit, ur domain.UserRepository, ar domain.AccountRepository, lr domain.LedgerRepository, sysEmail string) *PayoutService {
	return &PayoutService{
		atomicUnit:  atomicUnit,
		userRepo:    ur,
		accountRepo: ar,
		ledgerRepo:  lr,
		sysEmail:    sysEmail,
	}
}

func (ps *PayoutService) ExecutePayout(ctx context.Context, userID uuid.UUID, currency domain.Currency, amount int64, accountNum, bankCode string) (*domain.Transaction, error) {
	var payoutTxn *domain.Transaction

	err := ps.atomicUnit.Do(ctx, func(ctxWithTx context.Context) error {
		userAcct, err := getOrCreateAccountUtil(ps.accountRepo, ctxWithTx, userID, currency)
		if err != nil {
			return err
		}

		currentBalance, err := ps.ledgerRepo.GetAccountBalance(ctxWithTx, userAcct)
		if err != nil {
			return err
		}
		if currentBalance < amount {
			return ErrInsufficientFunds
		}

		sysUser, err := ps.userRepo.FindByEmail(ctxWithTx, ps.sysEmail)
		if err != nil {
			return err
		}
		systemAcct, err := getOrCreateAccountUtil(ps.accountRepo, ctxWithTx, sysUser.ID, currency)
		if err != nil {
			return err
		}

		// pending txn
		builder := domain.NewTransactionBuilder(domain.PAYOUT, domain.PENDING, uuid.New().String())
		builder.AddEntry(userAcct.ID, amount, domain.DEBIT, currency)
		builder.AddEntry(systemAcct.ID, amount, domain.CREDIT, currency)

		payoutTxn, err = builder.Build()
		if err != nil {
			return err
		}

		return ps.ledgerRepo.AppendTransaction(ctxWithTx, payoutTxn)
	})

	if err != nil {
		return nil, err
	}

	// async worker
	go ps.simulateBank(context.Background(), payoutTxn.ID, userID, amount, currency)

	return payoutTxn, nil
}

func (ps *PayoutService) simulateBank(ctx context.Context, txnID, userID uuid.UUID, amount int64, currency domain.Currency) {
	time.Sleep(5 * time.Second)

	ps.ledgerRepo.UpdateTransactionStatus(ctx, txnID, domain.PROCESSING)
	log.Info().Msgf("Payout %s processing", txnID)
	time.Sleep(2 * time.Second)

	success := rand.Float32() < 0.80

	if success {
		log.Info().Msgf("Payout %s succeeded", txnID)
		if err := ps.ledgerRepo.UpdateTransactionStatus(ctx, txnID, domain.SUCCESS); err != nil {
			log.Info().Msgf("Payout %s: Failed to update transaction status", txnID)
		}
	} else {
		log.Info().Msgf("Payout %s failed. Initiating reversal.", txnID)
		if err := ps.ledgerRepo.UpdateTransactionStatus(ctx, txnID, domain.FAILED); err != nil {
			log.Info().Msgf("Payout %s: Failed to update transaction status", txnID)
		}

		ps.executeReversal(ctx, txnID, userID, amount, currency)
	}
}

func (ps *PayoutService) executeReversal(ctx context.Context, originalTxnID, userID uuid.UUID, amount int64, currency domain.Currency) {
	err := ps.atomicUnit.Do(ctx, func(ctxWithTx context.Context) error {
		userAcct, err := getOrCreateAccountUtil(ps.accountRepo, ctxWithTx, userID, currency)
		if err != nil {
			return fmt.Errorf("failed to get/create user account: %w", err)
		}
		sysUser, err := ps.userRepo.FindByEmail(ctxWithTx, ps.sysEmail)
		if err != nil {
			return fmt.Errorf("failed to find system user: %w", err)
		}
		systemAcct, err := getOrCreateAccountUtil(ps.accountRepo, ctxWithTx, sysUser.ID, currency)
		if err != nil {
			return fmt.Errorf("failed to get/create system account: %w", err)
		}

		// reversal linked to via original reference
		reference := fmt.Sprintf("REVERSAL-%s", originalTxnID.String())
		builder := domain.NewTransactionBuilder(domain.REVERSAL, domain.SUCCESS, reference)

		builder.AddEntry(systemAcct.ID, amount, domain.DEBIT, currency)
		builder.AddEntry(userAcct.ID, amount, domain.CREDIT, currency)

		txn, err := builder.Build()
		if err != nil {
			return fmt.Errorf("failed to build reversal transaction: %w", err)
		}
		return ps.ledgerRepo.AppendTransaction(ctxWithTx, txn)
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("original_txn_id", originalTxnID.String()).
			Str("user_id", userID.String()).
			Int64("amount", amount).
			Str("currency", currency.String()).
			Msg("CRITICAL: Failed to execute ledger reversal. User funds stranded.")
	} else {
		log.Info().
			Str("original_txn_id", originalTxnID.String()).
			Msg("Reversal executed successfully")
	}
}
