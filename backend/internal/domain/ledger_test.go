package domain_test

import (
	"testing"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTransactionBuilder_BalancesCorrectly(t *testing.T) {
	acct1 := uuid.New()
	acct2 := uuid.New()

	tests := []struct {
		name      string
		txType    domain.TxnType
		status    domain.Status
		reference string
		entries   []struct {
			accountID uuid.UUID
			amount    int64
			entryType domain.Direction
			currency  domain.Currency
		}
		wantErr        error
		wantEntryCount int
	}{
		{
			name:      "Valid balanced transaction succeeds",
			txType:    domain.DEPOSIT,
			status:    domain.SUCCESS,
			reference: "ref-1",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{
				{acct1, 1000, domain.CREDIT, domain.USD},
				{acct2, 1000, domain.DEBIT, domain.USD},
			},
			wantErr:        nil,
			wantEntryCount: 2,
		},
		{
			name:      "Multi-currency valid transaction succeeds (FX)",
			txType:    domain.CONVERSION,
			status:    domain.SUCCESS,
			reference: "ref-2",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{
				{acct1, 100, domain.DEBIT, domain.USD},
				{acct2, 100, domain.CREDIT, domain.USD},
				{acct2, 150000, domain.DEBIT, domain.NGN},
				{acct1, 150000, domain.CREDIT, domain.NGN},
			},
			wantErr:        nil,
			wantEntryCount: 4,
		},
		{
			name:      "Unbalanced transaction fails",
			txType:    domain.DEPOSIT,
			status:    domain.SUCCESS,
			reference: "ref-3",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{
				{acct1, 1000, domain.CREDIT, domain.USD},
				{acct2, 900, domain.DEBIT, domain.USD},
			},
			wantErr:        domain.ErrTransactionImbalance,
			wantEntryCount: 0,
		},
		{
			name:      "Multi-currency unbalanced transaction fails (FX)",
			txType:    domain.CONVERSION,
			status:    domain.SUCCESS,
			reference: "ref-4",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{
				{acct1, 100, domain.DEBIT, domain.USD},
				{acct2, 100, domain.CREDIT, domain.USD},
				{acct2, 150000, domain.DEBIT, domain.NGN},
				{acct1, 150001, domain.CREDIT, domain.NGN},
			},
			wantErr:        domain.ErrTransactionImbalance,
			wantEntryCount: 0,
		},
		{
			name:      "Negative amounts fail",
			txType:    domain.DEPOSIT,
			status:    domain.SUCCESS,
			reference: "ref-5",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{
				{acct1, -1000, domain.CREDIT, domain.USD},
				{acct2, -1000, domain.DEBIT, domain.USD},
			},
			wantErr:        domain.ErrInvalidAmount,
			wantEntryCount: 0,
		},
		{
			name:      "Zero amounts fail",
			txType:    domain.PAYOUT,
			status:    domain.SUCCESS,
			reference: "ref-6",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{
				{acct1, 0, domain.CREDIT, domain.USD},
				{acct2, 0, domain.DEBIT, domain.USD},
			},
			wantErr:        domain.ErrInvalidAmount,
			wantEntryCount: 0,
		},
		{
			name:      "Transaction with 1 entry fails",
			txType:    domain.DEPOSIT,
			status:    domain.SUCCESS,
			reference: "ref-7",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{
				{acct1, 500, domain.CREDIT, domain.USD},
			},
			wantErr:        domain.ErrEmptyTransaction,
			wantEntryCount: 0,
		},
		{
			name:      "Transaction with 0 entries fails",
			txType:    domain.DEPOSIT,
			status:    domain.SUCCESS,
			reference: "ref-8",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{},
			wantErr:        domain.ErrEmptyTransaction,
			wantEntryCount: 0,
		},
		{
			name:      "Transaction with same side entries 1",
			txType:    domain.DEPOSIT,
			status:    domain.SUCCESS,
			reference: "ref-9",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{
				{acct1, 500, domain.CREDIT, domain.USD},
				{acct1, 200, domain.CREDIT, domain.USD},
				{acct1, 500, domain.CREDIT, domain.USD},
			},
			wantErr:        domain.ErrTransactionImbalance,
			wantEntryCount: 0,
		},
		{
			name:      "Transaction with same side entries 2",
			txType:    domain.DEPOSIT,
			status:    domain.SUCCESS,
			reference: "ref-10",
			entries: []struct {
				accountID uuid.UUID
				amount    int64
				entryType domain.Direction
				currency  domain.Currency
			}{
				{acct1, 500, domain.DEBIT, domain.USD},
				{acct1, 500, domain.DEBIT, domain.USD},
			},
			wantErr:        domain.ErrTransactionImbalance,
			wantEntryCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := domain.NewTransactionBuilder(tt.txType, tt.status, tt.reference)
			for _, e := range tt.entries {
				builder.AddEntry(e.accountID, e.amount, e.entryType, e.currency)
			}
			txn, err := builder.Build()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, txn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, txn)
				assert.Len(t, txn.Entries, tt.wantEntryCount)
			}
		})
	}
}
