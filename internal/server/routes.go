package server

import (
	"net/http"

	"github.com/Abdelrahiim/lms/internal/middleware"
)

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
