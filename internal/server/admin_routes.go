package server

import (
	"net/http"

	"github.com/Abdelrahiim/lms/internal/middleware"
)

// registerAdminRoutes handles system administration
func (s *Server) registerAdminRoutes(mux *http.ServeMux, globalMiddleware []middleware.Middleware) {
	// TODO: Initialize admin handler when implemented
	// adminHandler := handler.NewAdminHandler(s.db, s.queries, s.config)

	// User management
	// mux.HandleFunc("GET /api/v1/admin/users", chain(
	//     adminHandler.ListUsers,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireRole("admin"))...,
	// ))
	// mux.HandleFunc("PUT /api/v1/admin/users/{id}/role", chain(
	//     adminHandler.UpdateUserRole,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireRole("admin"), middleware.ValidateJSON[handler.UpdateRoleRequest])...,
	// ))
	// mux.HandleFunc("DELETE /api/v1/admin/users/{id}", chain(
	//     adminHandler.DeleteUser,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireRole("admin"))...,
	// ))

	// System settings
	// mux.HandleFunc("GET /api/v1/admin/settings", chain(
	//     adminHandler.GetSystemSettings,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireRole("admin"))...,
	// ))
	// mux.HandleFunc("PUT /api/v1/admin/settings", chain(
	//     adminHandler.UpdateSystemSettings,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireRole("admin"), middleware.ValidateJSON[handler.UpdateSettingsRequest])...,
	// ))
}
