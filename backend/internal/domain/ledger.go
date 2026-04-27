package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyTransaction     = errors.New("A transaction must contain at least 2 entries")
	ErrInvalidAmount        = errors.New("amount must be positive")
	ErrCurrencyMismatch     = errors.New("ledger entries within a transaction must share the same currency")
	ErrTransactionImbalance = errors.New("critical: ledger entries in tis transaction are not balanced")
)

type Account struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Currency  Currency
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LedgerEntry struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	TxnId     uuid.UUID
	Amount    int64
	Direction Direction
	Currency  Currency
	CreatedAt time.Time
}

type Transaction struct {
	ID        uuid.UUID
	Entries   []LedgerEntry
	Type      TxnType
	Status    Status
	Reference string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Ledger struct {
	ledgerEntries []LedgerEntry
}

type TransactionBuilder struct {
	txnType   TxnType
	status    Status
	reference string
	entries   []LedgerEntry
}

func NewTransactionBuilder(txnType TxnType, status Status, reference string) *TransactionBuilder {
	return &TransactionBuilder{
		txnType:   txnType,
		status:    status,
		reference: reference,
		entries:   make([]LedgerEntry, 0),
	}
}

func (tb *TransactionBuilder) AddEntry(accountID uuid.UUID, amount int64, direction Direction, currency Currency) {
	entry := LedgerEntry{
		AccountID: accountID,
		Amount:    amount,
		Direction: direction,
		Currency:  currency,
	}
	tb.entries = append(tb.entries, entry)
}

func (tb *TransactionBuilder) IsValidAmount() bool {
	for _, entry := range tb.entries {
		if entry.Amount <= 0 {
			return false
		}
	}
	return true
}

func (tb *TransactionBuilder) IsBalanced() bool {
	balances := make(map[Currency]int64)

	for _, entry := range tb.entries {
		if entry.Direction == CREDIT {
			balances[entry.Currency] += entry.Amount
		} else {
			balances[entry.Currency] -= entry.Amount
		}
	}

	for _, netBal := range balances {
		if netBal != 0 {
			return false
		}
	}
	return true
}

func (tb *TransactionBuilder) Build() (*Transaction, error) {
	if len(tb.entries) < 2 {
		return nil, ErrEmptyTransaction
	}

	if !tb.IsValidAmount() {
		return nil, ErrInvalidAmount
	}

	if !tb.IsBalanced() {
		return nil, ErrTransactionImbalance
	}

	return &Transaction{
		Entries:   tb.entries,
		Type:      tb.txnType,
		Status:    tb.status,
		Reference: tb.reference,
	}, nil
}
