package db

import (
	"context"
	"fmt"
	"time"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresLedgerRepo struct {
	db *sqlx.DB
}

type dbBalance struct {
	Currency string `db:"currency"`
	Balance  int64  `db:"balance"`
}

func NewPostgresLedgerRepo(db *sqlx.DB) *PostgresLedgerRepo {
	return &PostgresLedgerRepo{
		db: db,
	}
}

func (r *PostgresLedgerRepo) AppendTransaction(ctx context.Context, txn *domain.Transaction) error {
	q := getQuerier(ctx, r.db)
	now := time.Now().UTC()
	txnQuery := `
		INSERT INTO transactions(type, status, reference, created_at, updated_at)
		VALUES($1, $2, $3, $4, $4)
		RETURNING id;
	`

	// insert txn
	if err := q.QueryRowContext(ctx, txnQuery, txn.Type, txn.Status, txn.Reference, now).Scan(&txn.ID); err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	entryQuery := `
		INSERT INTO ledger_entries(account_id, amount, direction, transaction_id, currency, created_at)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	for _, entry := range txn.Entries {
		if err := q.QueryRowContext(ctx, entryQuery, entry.AccountID, entry.Amount, entry.Direction.String(), txn.ID, entry.Currency.String(), now).Scan(&entry.ID); err != nil {
			return fmt.Errorf("failed to insert ledger entry %w", err)
		}
		entry.TxnId = txn.ID
		entry.CreatedAt = now
	}

	txn.CreatedAt = now
	txn.UpdatedAt = now

	return nil
}

func (r *PostgresLedgerRepo) GetAccountBalance(ctx context.Context, account *domain.Account) (int64, error) {
	q := getQuerier(ctx, r.db)

	query := `
		SELECT COALESCE(SUM(
			CASE 
                WHEN direction = 'CREDIT' THEN amount 
                ELSE -amount 
            END
		), 0) 
		FROM ledger_entries
		WHERE account_id=$1
	`

	var bal int64
	if err := q.GetContext(ctx, &bal, query, account.ID); err != nil {
		return 0, fmt.Errorf("failed to calculate balance: %w", err)
	}

	return bal, nil
}

func (r *PostgresLedgerRepo) GetAllAccountBalances(ctx context.Context, userID uuid.UUID) (map[domain.Currency]int64, error) {
	q := getQuerier(ctx, r.db)

	query := `
        SELECT a.currency, COALESCE(SUM(
            CASE 
                WHEN le.direction = 'CREDIT' THEN le.amount 
                ELSE -le.amount 
            END
        ), 0) AS balance
        FROM accounts a
        LEFT JOIN ledger_entries le ON a.id = le.account_id
        WHERE a.user_id = $1
        GROUP BY a.currency;
    `

	var dbBalances []dbBalance
	if err := q.SelectContext(ctx, &dbBalances, query, userID); err != nil {
		return nil, fmt.Errorf("failed to query all balances: %w", err)
	}

	balances := make(map[domain.Currency]int64)
	for _, dbb := range dbBalances {
		if currency, ok := domain.GetCurrency(dbb.Currency); ok {
			balances[currency] = dbb.Balance
		}
	}

	return balances, nil
}
