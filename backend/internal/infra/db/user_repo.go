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

type PostgresUserRepo struct {
	db *sqlx.DB
}

func NewPostgresUserRepo(db *sqlx.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

type dbUser struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Firstname    string    `db:"first_name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (r *PostgresUserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users(email, password_hash, first_name, created_at) 
		VALUES($1, $2, $3, $4) 
		RETURNING id;
	`

	if err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.FirstName, time.Now().UTC()).Scan(&user.ID); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser dbUser

	query := `
		SELECT id, email, password_hash, first_name, created_at, created_at
		FROM users WHERE email=$1;
	`
	if err := r.db.GetContext(ctx, &dbUser, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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
