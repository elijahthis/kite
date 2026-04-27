package interfaces

import "net/http"

type Handlers struct {
	Auth    *AuthHandler
	Deposit *DepositHandler
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

	// handler := corsMiddleware(r.Router)

	// return handler
	return r.Router
}
