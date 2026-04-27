package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresQuoteRepo struct {
	db *sqlx.DB
}

func NewPostgresQuoteRepo(db *sqlx.DB) *PostgresQuoteRepo {
	return &PostgresQuoteRepo{db: db}
}

func (r *PostgresQuoteRepo) Save(ctx context.Context, quote *domain.Quote) error {
	query := `
        INSERT INTO quotes (id, user_id, source_currency, target_currency, provider_rate, exchange_rate, amount_in, amount_out, expires_at, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	_, err := r.db.ExecContext(ctx, query,
		quote.ID, quote.UserID, quote.SourceCurrency.String(), quote.TargetCurrency.String(),
		quote.ProviderRate, quote.ExchangeRate, quote.AmountIn, quote.AmountOut,
		quote.ExpiresAt, quote.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save quote: %w", err)
	}
	return nil
}

func (r *PostgresQuoteRepo) Get(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
	query := `
        SELECT id, user_id, source_currency, target_currency, provider_rate, exchange_rate, amount_in, amount_out, expires_at, created_at
        FROM quotes WHERE id = $1
    `

	var q domain.Quote
	var src, tgt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&q.ID, &q.UserID, &src, &tgt,
		&q.ProviderRate, &q.ExchangeRate, &q.AmountIn, &q.AmountOut,
		&q.ExpiresAt, &q.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrQuoteNotFound
		}
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	source, ok := domain.GetCurrency(src)
	if !ok {
		return nil, fmt.Errorf("invalid currency:")
	}
	target, ok := domain.GetCurrency(tgt)
	if !ok {
		return nil, fmt.Errorf("invalid currency:")
	}

	q.SourceCurrency = domain.Currency(source)
	q.TargetCurrency = domain.Currency(target)

	return &q, nil
}
