package application_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/elijahthis/kite/internal/application"
	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// THREAD-SAFE MOCKS
// ==========================================

// bypasses DB transactions
type mockAtomicUnit struct{}

func (m *mockAtomicUnit) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockLedgerRepo struct {
	mu           sync.Mutex
	txns         map[string]*domain.Transaction
	balances     map[uuid.UUID]int64
	ReversalSeen bool
	statusLog    []domain.Status
}

func newMockLedgerRepo() *mockLedgerRepo {
	return &mockLedgerRepo{
		txns:      make(map[string]*domain.Transaction),
		balances:  make(map[uuid.UUID]int64),
		statusLog: []domain.Status{},
	}
}

func (m *mockLedgerRepo) AppendTransaction(ctx context.Context, txn *domain.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// UNIQUE constraint
	if _, exists := m.txns[txn.Reference]; exists {
		return errors.New("unique constraint violation")
	}

	if txn.Type == domain.REVERSAL {
		m.ReversalSeen = true
	}

	m.txns[txn.Reference] = txn
	return nil
}

func (m *mockLedgerRepo) GetAccountBalance(ctx context.Context, account *domain.Account) (int64, error) {
	// plenty money
	return 1000000, nil
}
func (m *mockLedgerRepo) GetAllAccountBalances(ctx context.Context, userID uuid.UUID) (map[domain.Currency]int64, error) {
	return nil, nil
}
func (m *mockLedgerRepo) UpdateTransactionStatus(ctx context.Context, txnID uuid.UUID, status domain.Status) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.statusLog = append(m.statusLog, status)
	return nil
}
func (m *mockLedgerRepo) GetTransactionHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.TransactionHistory, error) {
	return nil, nil
}

type mockAccountRepo struct{}

func (m *mockAccountRepo) Create(ctx context.Context, account *domain.Account) error { return nil }
func (m *mockAccountRepo) GetByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency domain.Currency) (*domain.Account, error) {
	return &domain.Account{ID: uuid.New(), UserID: userID, Currency: currency}, nil
}
func (m *mockAccountRepo) LockAccount(ctx context.Context, accountID uuid.UUID) error { return nil }

type mockQuoteRepo struct {
	quote *domain.Quote
}

func (m *mockQuoteRepo) Save(ctx context.Context, quote *domain.Quote) error { return nil }
func (m *mockQuoteRepo) Get(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
	return m.quote, nil
}

type mockUserRepo struct{}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error { return nil }
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return &domain.User{}, nil
}

// TESTS
// ==========================================

func TestIdempotency_Works(t *testing.T) {
	mockLedger := newMockLedgerRepo()
	depositService := application.NewDepositService(
		&mockAccountRepo{}, mockLedger, &mockUserRepo{}, &mockAtomicUnit{}, "sys@grey.co",
	)

	userID := uuid.New()
	reference := "idempotent-webhook-123"

	// attempt 1
	err := depositService.ExecuteDeposit(context.Background(), userID, domain.USD, 5000, reference)
	assert.NoError(t, err)

	// attempt 2
	err = depositService.ExecuteDeposit(context.Background(), userID, domain.USD, 5000, reference)
	assert.ErrorIs(t, err, domain.ErrDuplicateReference)

	// 1 transaction only
	assert.Len(t, mockLedger.txns, 1)
}

func TestConcurrent_Conversion(t *testing.T) {
	mockLedger := newMockLedgerRepo()
	quoteID := uuid.New()
	userID := uuid.New()

	mockQuotes := &mockQuoteRepo{
		quote: &domain.Quote{
			ID:             quoteID, // txn ref
			UserID:         userID,
			SourceCurrency: domain.USD,
			TargetCurrency: domain.NGN,
			AmountIn:       100,
			AmountOut:      150000,
			ExpiresAt:      time.Now().Add(5 * time.Minute),
		},
	}

	conversionService := application.NewConversionService(
		nil, mockQuotes, &mockAtomicUnit{}, &mockUserRepo{}, &mockAccountRepo{}, mockLedger, "sys@grey.co",
	)

	// 500 concurrent requests, same quote
	var wg sync.WaitGroup
	successCount := 0
	duplicateCount := 0
	var mu sync.Mutex

	const TOTAL_REQ = 500

	for range TOTAL_REQ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := conversionService.ExecuteQuote(context.Background(), userID, quoteID)

			mu.Lock()
			if err == nil {
				successCount++
			} else if errors.Is(err, domain.ErrDuplicateReference) {
				duplicateCount++
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	assert.Equal(t, 1, successCount, "Only one concurrent execution should succeed")
	assert.Equal(t, TOTAL_REQ-1, duplicateCount, "All other attempts must return ErrDuplicateReference")
	assert.Len(t, mockLedger.txns, 1, "Only one conversion transaction should be written to the ledger")
}

func TestExpired_Quote(t *testing.T) {
	// quote that expired 1 minute ago
	userID := uuid.New()
	mockQuotes := &mockQuoteRepo{
		quote: &domain.Quote{
			ID:        uuid.New(),
			UserID:    userID,
			ExpiresAt: time.Now().Add(-1 * time.Minute),
		},
	}

	conversionService := application.NewConversionService(
		nil, mockQuotes, &mockAtomicUnit{}, &mockUserRepo{}, &mockAccountRepo{}, newMockLedgerRepo(), "sys@grey.co",
	)

	err := conversionService.ExecuteQuote(context.Background(), userID, mockQuotes.quote.ID)
	assert.ErrorIs(t, err, domain.ErrQuoteExpired)
}
func TestMismatchedUsers_Quote(t *testing.T) {
	// quote that expired 1 minute ago
	mockQuotes := &mockQuoteRepo{
		quote: &domain.Quote{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			ExpiresAt: time.Now().Add(-1 * time.Minute),
		},
	}

	conversionService := application.NewConversionService(
		nil, mockQuotes, &mockAtomicUnit{}, &mockUserRepo{}, &mockAccountRepo{}, newMockLedgerRepo(), "sys@grey.co",
	)

	err := conversionService.ExecuteQuote(context.Background(), uuid.New(), mockQuotes.quote.ID)
	assert.ErrorIs(t, err, domain.ErrQuoteBelongsToAnotherUser)
}

func TestFailed_Payout_Reversal_Unit(t *testing.T) {
	mockLedger := newMockLedgerRepo()
	payoutService := application.NewPayoutService(
		&mockAtomicUnit{}, &mockUserRepo{}, &mockAccountRepo{}, mockLedger, "sys@grey.co",
	)

	userID := uuid.New()
	originalTxnID := uuid.New()

	payoutService.ExecuteReversal(context.Background(), originalTxnID, userID, 5000, domain.NGN)

	// Assert that a new transaction of type REVERSAL was successfully appended
	assert.True(t, mockLedger.ReversalSeen, "A ledger reversal transaction should have been appended")

	// Assert the reference is correctly linked
	expectedReference := "REVERSAL-" + originalTxnID.String()
	_, exists := mockLedger.txns[expectedReference]
	assert.True(t, exists, "Reversal must be linked to the original transaction ID")
}

func TestFailed_Payout_Reversal_Integration(t *testing.T) {
	mockLedger := newMockLedgerRepo()

	payoutService := application.NewPayoutService(
		&mockAtomicUnit{}, &mockUserRepo{}, &mockAccountRepo{}, mockLedger, "sys@grey.co",
		application.WithBankSimulation(func() bool { return false }, 0),
	)

	userID := uuid.New()
	AMOUNT := 5000
	CURRENCY := domain.NGN
	txn, err := payoutService.ExecutePayout(
		context.Background(), userID, CURRENCY, int64(AMOUNT), "0123456789", "058",
	)
	assert.NoError(t, err, "payout initiation should succeed")
	assert.NotNil(t, txn)

	time.Sleep(20 * time.Millisecond)

	assert.Len(t, mockLedger.txns, 2, "ledger must contain the original payout + the reversal")
	assert.True(t, mockLedger.ReversalSeen, "a REVERSAL entry must be written after failure")

	// pending → PROCESSING → FAILED
	assert.Contains(t, mockLedger.statusLog, domain.PROCESSING)
	assert.Contains(t, mockLedger.statusLog, domain.FAILED)

	expectedRef := "REVERSAL-" + txn.ID.String()
	_, exists := mockLedger.txns[expectedRef]
	assert.True(t, exists, "reversal reference must be traceable to the original payout ID")
}

func TestSuccessful_Payout_No_Reversal(t *testing.T) {
	mockLedger := newMockLedgerRepo()

	payoutService := application.NewPayoutService(
		&mockAtomicUnit{}, &mockUserRepo{}, &mockAccountRepo{}, mockLedger, "sys@grey.co",
		application.WithBankSimulation(func() bool { return true }, 0),
	)

	_, err := payoutService.ExecutePayout(
		context.Background(), uuid.New(), domain.NGN, 5000, "0123456789", "058",
	)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	assert.False(t, mockLedger.ReversalSeen, "a successful payout must never write a reversal")
	assert.Len(t, mockLedger.txns, 1, "ledger must contain only the original payout")
	assert.Contains(t, mockLedger.statusLog, domain.PROCESSING)
	assert.Contains(t, mockLedger.statusLog, domain.SUCCESS)
}
