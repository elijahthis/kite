package db

import (
	"context"
	"database/sql"
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
	now := time.Now().UTC()
	txnQuery := `
		INSERT INTO transactions(type, status, reference, created_at, updated_at)
		VALUES($1, $2, $3, $4, $4)
		RETURNING id;
	`

	// begin txn
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// insert tx
	if err := tx.QueryRowContext(ctx, txnQuery, txn.Type, txn.Status, txn.Reference, now).Scan(&txn.ID); err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	entryQuery := `
		INSERT INTO ledger_entries(account_id, amount, direction, transaction_id, currency, created_at)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	for _, entry := range txn.Entries {
		if err := tx.QueryRowContext(ctx, entryQuery, entry.AccountID, entry.Amount, entry.Direction.String(), txn.ID, entry.Currency.String(), now).Scan(&entry.ID); err != nil {
			return fmt.Errorf("failed to insert ledger entry", err)
		}
		entry.TxnId = txn.ID
		entry.CreatedAt = now
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	txn.CreatedAt = now
	txn.UpdatedAt = now

	return nil
}

func (r *PostgresLedgerRepo) GetAccountBalance(ctx context.Context, account *domain.Account) (int64, error) {
	return 0, nil
}
