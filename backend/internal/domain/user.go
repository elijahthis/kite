package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists  = errors.New("a user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrDuplicateReference = errors.New("a transaction with this reference already exists")
)

type User struct {
	ID           uuid.UUID `db:"id"`
	FirstName    string    `db:"first_name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type AuthUseCase interface {
	Register(ctx context.Context, email string, plainPassword string, firstName string) (*User, error)
	Login(ctx context.Context, email string, plainPassword string) (string, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type TokenGenerator interface {
	GenerateToken(userID uuid.UUID) (string, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}
