package domain

import (
	"context"

	"github.com/google/uuid"
)

type AtomicUnit interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	GetByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency Currency) (*Account, error)
}
type LedgerRepository interface {
	AppendTransaction(ctx context.Context, txn *Transaction) error
	GetAccountBalance(ctx context.Context, account *Account) (int64, error)
	GetAllAccountBalances(ctx context.Context, userID uuid.UUID) (map[Currency]int64, error)
}

type QuoteRepository interface {
	Save(ctx context.Context, quote *Quote) error
	Get(ctx context.Context, id uuid.UUID) (*Quote, error)
}
