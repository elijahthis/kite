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
)

type User struct {
	ID           uuid.UUID
	FirstName    string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
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
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}
