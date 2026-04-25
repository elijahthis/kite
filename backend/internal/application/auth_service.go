package application

import (
	"context"

	"github.com/elijahthis/kite/internal/domain"
)

type AuthService struct {
	hasher         domain.PasswordHasher
	tokenGenerator domain.TokenGenerator
	repo           domain.UserRepository
}

func NewAuthService(h domain.PasswordHasher, t domain.TokenGenerator, r domain.UserRepository) *AuthService {
	return &AuthService{
		hasher:         h,
		tokenGenerator: t,
		repo:           r,
	}
}

func (as *AuthService) Register(ctx context.Context, email string, plainPassword string, firstName string) (*domain.User, error) {
	existingUser, err := as.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	hashedPassword, err := as.hasher.Hash(plainPassword)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		FirstName:    firstName,
		Email:        email,
		PasswordHash: hashedPassword,
	}
	if err := as.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (as *AuthService) Login(ctx context.Context, email string, plainPassword string) (string, error) {
	user, err := as.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", domain.ErrInvalidCredentials
	}
	if !as.hasher.Compare(user.PasswordHash, plainPassword) {
		return "", domain.ErrInvalidCredentials
	}

	return as.tokenGenerator.GenerateToken(user.ID)
}
