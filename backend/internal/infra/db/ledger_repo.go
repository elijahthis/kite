package db

import (
	"context"
	"fmt"
	"time"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresLedgerRepo struct {
	db *sqlx.DB
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
