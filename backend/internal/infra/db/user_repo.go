package db

import (
	"context"
	"fmt"
	"time"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresUserRepo struct {
	db *sqlx.DB
}

type dbUser struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Firstname    string    `db:"firstname"`
	CreatedAt    time.Time `db:"updated_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (r *PostgresUserRepo) Create(ctx context.Context, user *domain.User) error {
	var id string
	query := `
		INSERT INTO users(email, password_hash, firstname, created_at) 
		VALUES($1, $2, $3, $4) 
		RETURNING id;
	`

	if err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.FirstName, time.Now().UTC()).Scan(&id); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser dbUser

	query := `
		SELECT id, email, password_hash, created_at
		FROM users WHERE email=$1;
	`
	if err := r.db.GetContext(ctx, &dbUser, query, email); err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}
	return &domain.User{
		ID:           dbUser.ID,
		FirstName:    dbUser.Firstname,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}, nil
}
