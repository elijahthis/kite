package interfaces

import "net/http"

type Handlers struct {
	Auth       *AuthHandler
	Deposit    *DepositHandler
	Conversion *ConversionHandler
	Wallet     *WalletHandler
	Payout     *PayoutHandler
	History    *HistoryHandler
}

type Router struct {
	Router  *http.ServeMux
	Handler http.Handler
}

func NewRouter() *Router {
	mux := http.NewServeMux()

	return &Router{
		Router: mux,
	}
}

func (r *Router) SetupRouter(h Handlers, jwtSecret string) {
	if h.Auth != nil {
		r.Router.HandleFunc("POST /api/v1/auth/register", ApplyMiddleware(h.Auth.Register, RequestLogger))
		r.Router.HandleFunc("POST /api/v1/auth/login", ApplyMiddleware(h.Auth.Login, RequestLogger))
	}

	if h.Deposit != nil {
		r.Router.HandleFunc("POST /api/v1/deposits", ApplyMiddleware(h.Deposit.Create, RequireAuth(jwtSecret), RequestLogger))
	}

	if h.Conversion != nil {
		r.Router.HandleFunc("POST /api/v1/conversions/quote", ApplyMiddleware(h.Conversion.GenerateQuote, RequireAuth(jwtSecret), RequestLogger))
		r.Router.HandleFunc("POST /api/v1/conversions/execute", ApplyMiddleware(h.Conversion.ExecuteQuote, RequireAuth(jwtSecret), RequestLogger))
	}
	if h.Wallet != nil {
		r.Router.HandleFunc("GET /api/v1/balances", ApplyMiddleware(h.Wallet.GetBalances, RequireAuth(jwtSecret), RequestLogger))
	}
	if h.Payout != nil {
		r.Router.HandleFunc("POST /api/v1/payouts", ApplyMiddleware(h.Payout.Create, RequireAuth(jwtSecret), RequestLogger))
	}
	if h.History != nil {
		r.Router.HandleFunc("GET /api/v1/transactions", ApplyMiddleware(h.History.GetHistory, RequireAuth(jwtSecret), RequestLogger))
	}

	handler := corsMiddleware(r.Router)

	r.Handler = handler
	// return r.Router
}
