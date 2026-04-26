package config

import (
	"context"

	"github.com/elijahthis/kite/internal/application"
	"github.com/elijahthis/kite/internal/domain"
	"github.com/elijahthis/kite/internal/infra/crypto"
	"github.com/elijahthis/kite/internal/infra/db"
	interfaces "github.com/elijahthis/kite/internal/interfaces/kite_http"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Factory struct {
	Config         *Config
	DB             *sqlx.DB
	Repos          *Repos
	Hasher         domain.PasswordHasher
	TokenGenerator domain.TokenGenerator
	services       *application.Services
	handlers       *interfaces.Handlers
	Router         *interfaces.Router
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

	f := &Factory{
		Config:         cfg,
		DB:             dbx,
		Repos:          repos,
		Hasher:         hasher,
		TokenGenerator: tokenGen,
	}

	f.services = generateServices(f)
	f.handlers = generateHandlers(f)
	f.Router = newRouter(f)

	return f
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

func generateServices(f *Factory) *application.Services {
	return &application.Services{
		Auth: application.NewAuthService(f.Hasher, f.TokenGenerator, f.Repos.UserRepo),
	}
}

func generateHandlers(f *Factory) *interfaces.Handlers {
	return &interfaces.Handlers{
		Auth: interfaces.NewAuthHandler(*f.services.Auth),
	}
}

func newRouter(f *Factory) *interfaces.Router {
	r := interfaces.NewRouter()
	r.SetupRouter(*f.handlers)

	return r
}
