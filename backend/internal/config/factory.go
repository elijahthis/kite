package config

import (
	"context"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/elijahthis/kite/internal/infra/crypto"
	"github.com/elijahthis/kite/internal/infra/db"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Factory struct {
	config         *Config
	DB             *sqlx.DB
	Repos          *Repos
	Hasher         domain.PasswordHasher
	TokenGenerator domain.TokenGenerator
}

type Repos struct {
	UserRepo domain.UserRepository
}

func NewFactory() *Factory {
	cfg := LoadConfig()

	// infra
	dbx, err := db.NewDBPool(context.Background(), cfg.DB_CONN_STR)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize DB")
	}

	repos := newRepos(dbx)
	hasher := newHasher()
	tokenGen := newTokenGenerator(cfg.JWT_SECRET_KEY)

	return &Factory{
		config:         cfg,
		DB:             dbx,
		Repos:          repos,
		Hasher:         hasher,
		TokenGenerator: tokenGen,
	}
}

func newRepos(dbx *sqlx.DB) *Repos {
	userRepo := db.NewPostgresUserRepo(dbx)

	return &Repos{
		UserRepo: userRepo,
	}
}

func newHasher() domain.PasswordHasher {
	return crypto.NewBcryptHasher()
}

func newTokenGenerator(secretKey string) domain.TokenGenerator {
	return crypto.NewJWTTokenGenerator(secretKey)
}
