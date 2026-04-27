package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type contextKey string

const txKey contextKey = "kite_tx_key"


type SQLXAtomicUnit struct {
	db *sqlx.DB
}

func NewAtomicUnit(db *sqlx.DB) *SQLXAtomicUnit {
	return &SQLXAtomicUnit{db: db}
}

func (u *SQLXAtomicUnit) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	// Check for  transaction into the ctx
	if ctx.Value(txKey) != nil {
		return fn(ctx)
	}

	tx, err := u.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Add  transaction into the ctx
	ctxWithTx := context.WithValue(ctx, txKey, tx)
	err = fn(ctxWithTx)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback also failed: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Querier can either be transaction or DBB
type querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func getQuerier(ctx context.Context, db *sqlx.DB) querier {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
