package crypto

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTTokenGenerator struct {
	secretKey string
}

func NewJWTTokenGenerator(secretKey string) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		secretKey: secretKey,
	}
}

func (t *JWTTokenGenerator) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
