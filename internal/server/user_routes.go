package server

import (
	"net/http"

	"github.com/Abdelrahiim/lms/internal/middleware"
)

// registerUserRoutes handles user profile and settings
func (s *Server) registerUserRoutes(mux *http.ServeMux, globalMiddleware []middleware.Middleware) {
	// TODO: Initialize user handler when implemented
	// userHandler := handler.NewUserHandler(s.db, s.queries, s.config)

	// User profile endpoints
	// mux.HandleFunc("GET /api/v1/users/profile", chain(
	//     userHandler.GetProfile,
	//     append(globalMiddleware, middleware.RequireAuth)...,
	// ))
	// mux.HandleFunc("PUT /api/v1/users/profile", chain(
	//     userHandler.UpdateProfile,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.ValidateJSON[handler.UpdateProfileRequest])...,
	// ))
	// mux.HandleFunc("POST /api/v1/users/avatar", chain(
	//     userHandler.UploadAvatar,
	//     append(globalMiddleware, middleware.RequireAuth)...,
	// ))
}
