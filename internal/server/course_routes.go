package server

import (
	"net/http"

	"github.com/Abdelrahiim/lms/internal/middleware"
)

// registerCourseRoutes handles course management and enrollment
func (s *Server) registerCourseRoutes(mux *http.ServeMux, globalMiddleware []middleware.Middleware) {
	// TODO: Initialize course handler when implemented
	// courseHandler := handler.NewCourseHandler(s.db, s.queries, s.config)

	// Course discovery and enrollment
	// mux.HandleFunc("GET /api/v1/courses", chain(
	//     courseHandler.ListCourses,
	//     globalMiddleware...,
	// ))
	// mux.HandleFunc("GET /api/v1/courses/{id}", chain(
	//     courseHandler.GetCourse,
	//     globalMiddleware...,
	// ))
	// mux.HandleFunc("POST /api/v1/courses/{id}/enroll", chain(
	//     courseHandler.EnrollInCourse,
	//     append(globalMiddleware, middleware.RequireAuth)...,
	// ))

	// Course content (modules and lessons)
	// mux.HandleFunc("GET /api/v1/courses/{id}/modules", chain(
	//     courseHandler.GetCourseModules,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireEnrollment)...,
	// ))
	// mux.HandleFunc("GET /api/v1/modules/{id}/lessons", chain(
	//     courseHandler.GetModuleLessons,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireEnrollment)...,
	// ))
	// mux.HandleFunc("POST /api/v1/lessons/{id}/complete", chain(
	//     courseHandler.CompleteLesson,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireEnrollment)...,
	// ))

	// Instructor course management
	// mux.HandleFunc("POST /api/v1/courses", chain(
	//     courseHandler.CreateCourse,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireRole("instructor"), middleware.ValidateJSON[handler.CreateCourseRequest])...,
	// ))
	// mux.HandleFunc("PUT /api/v1/courses/{id}", chain(
	//     courseHandler.UpdateCourse,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireInstructor, middleware.ValidateJSON[handler.UpdateCourseRequest])...,
	// ))
}
