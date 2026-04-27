package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresAccountRepo struct {
	db *sqlx.DB
}

type dbAccount struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPostgresAccountRepo(db *sqlx.DB) *PostgresAccountRepo {
	return &PostgresAccountRepo{
		db: db,
	}
}

func (ar *PostgresAccountRepo) Create(ctx context.Context, account *domain.Account) error {
	now := time.Now().UTC()
	query := `
		INSERT INTO accounts(user_id, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $3)
		RETURNING id;
	`

	if err := ar.db.QueryRowContext(ctx, query, account.UserID, account.Currency.String(), now).Scan(&account.ID); err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	account.CreatedAt = now
	account.UpdatedAt = now

	return nil
}

func (ar *PostgresAccountRepo) GetByUserIdAndCurrency(ctx context.Context, userID uuid.UUID, currency domain.Currency) (*domain.Account, error) {
	var dbAccount dbAccount

	query := `
		SELECT id, user_id, currency, created_at, updated_at
		FROM accounts
		WHERE user_id=$1 AND currency=$2;
	`

	if err := ar.db.QueryRowContext(ctx, query, userID, currency.String()).Scan(&dbAccount.ID, &dbAccount.UserID, &dbAccount.Currency, &dbAccount.CreatedAt, &dbAccount.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}

	parsedCurrency, ok := domain.GetCurrency(dbAccount.Currency)
	if !ok {
		return nil, fmt.Errorf("invalid currency:")
	}

	return &domain.Account{
		ID:        dbAccount.ID,
		UserID:    dbAccount.UserID,
		Currency:  parsedCurrency,
		CreatedAt: dbAccount.CreatedAt,
		UpdatedAt: dbAccount.UpdatedAt,
	}, nil
}
