package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Abdelrahiim/lms/internal/config"
	"github.com/Abdelrahiim/lms/internal/database"
	_ "github.com/lib/pq"
)

type Server struct {
	config     *config.Config
	db         *sql.DB
	queries    *database.Queries
	httpServer *http.Server
}

func New(cfg *config.Config) (*Server, error) {
	// Open database connection
	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create SQLC queries
	queries := database.New(db)

	s := &Server{
		config:  cfg,
		db:      db,
		queries: queries,
	}

	// Setup routes
	mux := s.setupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return s, nil
}

func (s *Server) Start() error {
	// Graceful shutdown
	go s.handleShutdown()

	log.Printf("Server starting on port %s", s.config.Server.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) handleShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server shutdown complete")
}
