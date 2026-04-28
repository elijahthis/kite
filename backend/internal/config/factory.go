package config

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/elijahthis/kite/internal/application"
	"github.com/elijahthis/kite/internal/domain"
	"github.com/elijahthis/kite/internal/infra/crypto"
	"github.com/elijahthis/kite/internal/infra/db"
	"github.com/elijahthis/kite/internal/infra/fx"
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
	UserRepo    domain.UserRepository
	AccountRepo domain.AccountRepository
	LedgerRepo  domain.LedgerRepository
	AtomicUnit  domain.AtomicUnit
	FxProvider  domain.FXRateProvider
	QuoteRepo   domain.QuoteRepository
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
	accountRepo := db.NewPostgresAccountRepo(dbx)
	ledgerRepo := db.NewPostgresLedgerRepo(dbx)
	atomicUnit := db.NewAtomicUnit(dbx)

	fxProvider := fx.NewERAPIProvider(5 * time.Minute)
	quoteRepo := db.NewPostgresQuoteRepo(dbx)

	return &Repos{
		UserRepo:    userRepo,
		AccountRepo: accountRepo,
		LedgerRepo:  ledgerRepo,
		AtomicUnit:  atomicUnit,
		FxProvider:  fxProvider,
		QuoteRepo:   quoteRepo,
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
		Auth:       application.NewAuthService(f.Hasher, f.TokenGenerator, f.Repos.UserRepo),
		Deposit:    application.NewDepositService(f.Repos.AccountRepo, f.Repos.LedgerRepo, f.Repos.UserRepo, f.Repos.AtomicUnit, f.Config.SYSTEM_USER_EMAIL),
		Conversion: application.NewConversionService(f.Repos.FxProvider, f.Repos.QuoteRepo, f.Repos.AtomicUnit, f.Repos.UserRepo, f.Repos.AccountRepo, f.Repos.LedgerRepo, f.Config.SYSTEM_USER_EMAIL),
		Wallet:     application.NewWalletService(f.Repos.LedgerRepo),
		Payout:     application.NewPayoutService(f.Repos.AtomicUnit, f.Repos.UserRepo, f.Repos.AccountRepo, f.Repos.LedgerRepo, f.Config.SYSTEM_USER_EMAIL, application.WithBankSimulation(func() bool { return rand.Float32() < 0.5 }, 2*time.Second)),
		History:    application.NewHistoryService(f.Repos.LedgerRepo),
	}
}

func generateHandlers(f *Factory) *interfaces.Handlers {
	return &interfaces.Handlers{
		Auth:       interfaces.NewAuthHandler(*f.services.Auth),
		Deposit:    interfaces.NewDepositHandler(*f.services.Deposit),
		Conversion: interfaces.NewConversionHandler(*f.services.Conversion),
		Wallet:     interfaces.NewWalletHandler(*f.services.Wallet),
		Payout:     interfaces.NewPayoutHandler(*f.services.Payout),
		History:    interfaces.NewHistoryHandler(*f.services.History),
	}
}

func newRouter(f *Factory) *interfaces.Router {
	r := interfaces.NewRouter()
	r.SetupRouter(*f.handlers, f.Config.JWT_SECRET_KEY)

	return r
}
