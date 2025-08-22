package server

import (
	"net/http"

	"github.com/Abdelrahiim/lms/internal/handler"
	"github.com/Abdelrahiim/lms/internal/middleware"
)

// registerAuthRoutes handles authentication and session management
func (s *Server) registerAuthRoutes(mux *http.ServeMux, globalMiddleware []middleware.Middleware) {
	authHandler := handler.NewAuthHandler(s.db, s.queries, s.config)

	// Authentication endpoints
	mux.HandleFunc("POST /api/v1/auth/register", chain(
		authHandler.Register,
		append(globalMiddleware, middleware.ValidateJSON[handler.RegisterRequest])...,
	))

	mux.HandleFunc("POST /api/v1/auth/login", chain(
		authHandler.Login,
		append(globalMiddleware, middleware.ValidateJSON[handler.LoginRequest])...,
	))

	mux.HandleFunc("POST /api/v1/auth/logout", chain(
		authHandler.Logout,
		globalMiddleware...,
	))

	mux.HandleFunc("POST /api/v1/auth/refresh", chain(
		authHandler.Refresh,
		globalMiddleware...,
	))

	// Password management
	// mux.HandleFunc("POST /api/v1/auth/forgot-password", chain(
	//     authHandler.ForgotPassword,
	//     append(globalMiddleware, middleware.ValidateJSON[handler.ForgotPasswordRequest])...,
	// ))
	// mux.HandleFunc("POST /api/v1/auth/reset-password", chain(
	//     authHandler.ResetPassword,
	//     append(globalMiddleware, middleware.ValidateJSON[handler.ResetPasswordRequest])...,
	// ))
}
