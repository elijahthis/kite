package domain

import (
	"time"

	"github.com/google/uuid"
)

type AccountType int
type Direction int
type TxnType int
type Status int
type Currency int

const (
	CHECKING AccountType = iota
	ASSET
)
const (
	DEBIT Direction = iota
	CREDIT
)
const (
	DEPOSIT TxnType = iota
	PAYOUT
)
const (
	PENDING Status = iota
	PROCESSING
	SUCCESSFUL
	FAILED
)
const (
	USD Currency = iota
	GBP
	NGN
	KES
)

type Account struct {
	accountType AccountType
	id          uuid.UUID
	userID      uuid.UUID
	currency    Currency
	createdAt   time.Time
	UpdatedAt   time.Time
}

type LedgerEntry struct {
	accountID uuid.UUID
	amount    int
	direction Direction
	txId      uuid.UUID
	currency  Currency
	createdAt time.Time
}

type Transaction struct {
	id        uuid.UUID
	entries   []LedgerEntry
	txnType   TxnType
	Status    Status
	reference string
	createdAt time.Time
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

func (tb *TransactionBuilder) AddEntry(accountID uuid.UUID, amount int, direction Direction, currency Currency) {
	entry := LedgerEntry{
		accountID: accountID,
		amount:    amount,
		direction: direction,
		currency:  currency,
	}
	tb.entries = append(tb.entries, entry)
}

func (tb *TransactionBuilder) Build() (*Transaction, error) {
	return &Transaction{
		entries:   tb.entries,
		txnType:   tb.txnType,
		Status:    tb.status,
		reference: tb.reference,
	}, nil
}
