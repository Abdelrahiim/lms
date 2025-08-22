package server

import (
	"net/http"

	"github.com/Abdelrahiim/lms/internal/middleware"
)

// registerForumRoutes handles discussion forums and community features
func (s *Server) registerForumRoutes(mux *http.ServeMux, globalMiddleware []middleware.Middleware) {
	// TODO: Initialize forum handler when implemented
	// forumHandler := handler.NewForumHandler(s.db, s.queries, s.config)

	// Forum browsing and participation
	// mux.HandleFunc("GET /api/v1/courses/{id}/forum", chain(
	//     forumHandler.GetCourseForum,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireEnrollment)...,
	// ))
	// mux.HandleFunc("POST /api/v1/courses/{id}/threads", chain(
	//     forumHandler.CreateThread,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireEnrollment, middleware.ValidateJSON[handler.CreateThreadRequest])...,
	// ))
	// mux.HandleFunc("POST /api/v1/threads/{id}/posts", chain(
	//     forumHandler.CreatePost,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireEnrollment, middleware.ValidateJSON[handler.CreatePostRequest])...,
	// ))

	// Moderation endpoints
	// mux.HandleFunc("POST /api/v1/posts/{id}/report", chain(
	//     forumHandler.ReportPost,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.ValidateJSON[handler.ReportPostRequest])...,
	// ))
	// mux.HandleFunc("DELETE /api/v1/posts/{id}", chain(
	//     forumHandler.DeletePost,
	//     append(globalMiddleware, middleware.RequireAuth, middleware.RequireModerator)...,
	// ))
}
