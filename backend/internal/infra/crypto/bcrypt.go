package crypto

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{
		cost: bcrypt.DefaultCost,
	}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	h_bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(h_bytes), nil
}

func (h *BcryptHasher) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
