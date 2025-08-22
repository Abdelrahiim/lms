package server

import (
	"net/http"

	"github.com/Abdelrahiim/lms/internal/middleware"
)

// RegisterRoutes sets up all application routes with proper middleware chaining
func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Global middleware stack
	globalMiddleware := []middleware.Middleware{
		middleware.RequestID,
		middleware.Logger,
		middleware.Recovery,
		middleware.CORS,
	}

	// Health check endpoint
	mux.HandleFunc("/health", chain(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}, globalMiddleware...))

	// Register route groups
	s.registerAuthRoutes(mux, globalMiddleware)
	s.registerUserRoutes(mux, globalMiddleware)
	s.registerCourseRoutes(mux, globalMiddleware)
	s.registerAssessmentRoutes(mux, globalMiddleware)
	s.registerForumRoutes(mux, globalMiddleware)
	s.registerAnalyticsRoutes(mux, globalMiddleware)
	s.registerAdminRoutes(mux, globalMiddleware)

	return mux
}

// Helper function to chain middleware properly
func chain(f http.HandlerFunc, middlewares ...middleware.Middleware) http.HandlerFunc {
	return middleware.Chain(f, middlewares...)
}
