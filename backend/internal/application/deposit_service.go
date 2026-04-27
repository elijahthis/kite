package application

import (
	"context"
	"strings"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
)

type DepositService struct {
	accountRepo domain.AccountRepository
	ledgerRepo  domain.LedgerRepository
	userRepo    domain.UserRepository
	sysEmail    string
	atomicUnit  domain.AtomicUnit
}

func NewDepositService(ar domain.AccountRepository, lr domain.LedgerRepository, ur domain.UserRepository, at domain.AtomicUnit, sysEmail string) *DepositService {
	return &DepositService{
		accountRepo: ar,
		ledgerRepo:  lr,
		userRepo:    ur,
		atomicUnit:  at,
		sysEmail:    sysEmail,
	}
}

func (ds *DepositService) ExecuteDeposit(ctx context.Context, userId uuid.UUID, currency domain.Currency, amount int64, reference string) error {

	return ds.atomicUnit.Do(ctx, func(ctxWithTx context.Context) error {
		sysUser, err := ds.userRepo.FindByEmail(ctxWithTx, ds.sysEmail)
		if err != nil {
			return err
		}

		userAcct, err := ds.getOrCreateAccount(ctxWithTx, userId, currency)
		if err != nil {
			return err
		}

		systemAcct, err := ds.getOrCreateAccount(ctxWithTx, sysUser.ID, currency)
		if err != nil {
			return err
		}

		builder := domain.NewTransactionBuilder(domain.DEPOSIT, domain.SUCCESS, reference)
		builder.AddEntry(userAcct.ID, amount, domain.CREDIT, currency)
		builder.AddEntry(systemAcct.ID, amount, domain.DEBIT, currency)
		txn, err := builder.Build()
		if err != nil {
			return err
		}

		if err := ds.ledgerRepo.AppendTransaction(ctxWithTx, txn); err != nil {
			if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
				return domain.ErrDuplicateReference
			}
			return err
		}

		return nil
	})

}

func (ds *DepositService) getOrCreateAccount(ctx context.Context, ownerID uuid.UUID, currency domain.Currency) (*domain.Account, error) {
	return getOrCreateAccountUtil(ds.accountRepo, ctx, ownerID, currency)
}
