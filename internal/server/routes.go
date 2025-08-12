package server

import (
	"net/http"

	"github.com/Abdelrahiim/lms/internal/handler"
	"github.com/Abdelrahiim/lms/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Auth routes with validation
	authHandler := handler.NewAuthHandler(s.db, s.queries, s.config)
	
	// Apply validation middleware to specific routes
	mux.HandleFunc("POST /auth/register", 
		middleware.Chain(
			authHandler.Register,
			middleware.ValidateJSON[handler.RegisterRequest],
			middleware.Logger,
			middleware.RequestID,
		),
	)
	
	mux.HandleFunc("POST /auth/login", 
		middleware.Chain(
			authHandler.Login,
			middleware.ValidateJSON[handler.LoginRequest],
			middleware.Logger,
			middleware.RequestID,
		),
	)

	return mux
}

func (s *Server) setupRoutes() http.Handler {
	mux := http.NewServeMux()
	globalMiddleware := []middleware.Middleware{
		middleware.RequestID,
		middleware.Logger,
		middleware.Recovery,
		middleware.CORS,
	}

	// Initialize services
	mux.HandleFunc("/health", chain(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}, globalMiddleware...))

	return mux
}

// Helper function to chain middleware
func chain(f http.HandlerFunc, middlewares ...middleware.Middleware) http.HandlerFunc {
	return middleware.Chain(f, middlewares...)
}
