package interfaces

import "net/http"

type Handlers struct {
	Auth       *AuthHandler
	Deposit    *DepositHandler
	Conversion *ConversionHandler
	Wallet     *WalletHandler
	Payout     *PayoutHandler
}

type Router struct {
	Router *http.ServeMux
}

func NewRouter() *Router {
	mux := http.NewServeMux()

	return &Router{
		Router: mux,
	}
}

func (r *Router) SetupRouter(h Handlers, jwtSecret string) http.Handler {
	if h.Auth != nil {
		r.Router.HandleFunc("POST /api/v1/auth/register", h.Auth.Register)
		r.Router.HandleFunc("POST /api/v1/auth/login", h.Auth.Login)
	}

	if h.Deposit != nil {
		r.Router.HandleFunc("POST /api/v1/deposit", RequireAuth(jwtSecret)(h.Deposit.Create))
	}

	if h.Conversion != nil {
		r.Router.HandleFunc("POST /api/v1/conversions/quote", RequireAuth(jwtSecret)(h.Conversion.GenerateQuote))
		r.Router.HandleFunc("POST /api/v1/conversions/execute", RequireAuth(jwtSecret)(h.Conversion.ExecuteQuote))
	}
	if h.Wallet != nil {
		r.Router.HandleFunc("GET /api/v1/balances", RequireAuth(jwtSecret)(h.Wallet.GetBalances))
	}
	if h.Payout != nil {
		r.Router.HandleFunc("POST /api/v1/payouts", RequireAuth(jwtSecret)(h.Payout.Create))
	}

	handler := corsMiddleware(r.Router)

	return handler
	// return r.Router
}
