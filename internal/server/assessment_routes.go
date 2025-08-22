package server

import (
	"net/http"

	"github.com/Abdelrahiim/lms/internal/middleware"
)

// registerAssessmentRoutes handles quizzes, assignments, and grading
func (s *Server) registerAssessmentRoutes(mux *http.ServeMux, globalMiddleware []middleware.Middleware) {
	// TODO: Initialize assessment handler when implemented
	// assessmentHandler := handler.NewAssessmentHandler(s.db, s.queries, s.config)

	// Student assessment endpoints
	// mux.HandleFunc("GET /api/v1/courses/{id}/assessments", chain(
	//     assessmentHandler.ListAssessments,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireEnrollment)...,
	// ))
	// mux.HandleFunc("POST /api/v1/assessments/{id}/attempt", chain(
	//     assessmentHandler.StartAttempt,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireEnrollment)...,
	// ))
	// mux.HandleFunc("POST /api/v1/attempts/{id}/submit", chain(
	//     assessmentHandler.SubmitAttempt,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.ValidateJSON[handler.SubmitAttemptRequest])...,
	// ))

	// Instructor assessment management
	// mux.HandleFunc("POST /api/v1/courses/{id}/assessments", chain(
	//     assessmentHandler.CreateAssessment,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireInstructor, middleware.ValidateJSON[handler.CreateAssessmentRequest])...,
	// ))
	// mux.HandleFunc("GET /api/v1/assessments/{id}/results", chain(
	//     assessmentHandler.GetResults,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireInstructor)...,
	// ))
}
