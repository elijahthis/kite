package interfaces

import "net/http"

type Handlers struct {
	Auth *AuthHandler
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

func (r *Router) SetupRouter(h Handlers) http.Handler {
	if h.Auth != nil {
		r.Router.HandleFunc("POST /api/v1/auth/register", h.Auth.Register)
		r.Router.HandleFunc("POST /api/v1/auth/login", h.Auth.Login)
	}

	// handler := corsMiddleware(r.Router)

	// return handler
	return r.Router
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// will add logging middleware later
