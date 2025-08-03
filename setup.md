# Simple Standard Library HTTP Implementation with dotenv

## 1. Simple Configuration with dotenv

### internal/config/config.go
```go
package config

import (
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Auth     AuthConfig
    Storage  StorageConfig
}

type ServerConfig struct {
    Port         string
    Environment  string
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Database string
    SSLMode  string
}

type AuthConfig struct {
    JWTSecret          string
    JWTExpiry          time.Duration
    RefreshTokenExpiry time.Duration
    BcryptCost         int
}

type StorageConfig struct {
    UploadPath string
    MaxSize    int64
}

// Load loads configuration from .env file
func Load() (*Config, error) {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, using environment variables")
    }

    cfg := &Config{
        Server: ServerConfig{
            Port:         getEnv("PORT", "8080"),
            Environment:  getEnv("ENVIRONMENT", "development"),
            ReadTimeout:  getDurationEnv("READ_TIMEOUT", 15*time.Second),
            WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getIntEnv("DB_PORT", 5432),
            User:     mustGetEnv("DB_USER"),
            Password: mustGetEnv("DB_PASSWORD"),
            Database: mustGetEnv("DB_NAME"),
            SSLMode:  getEnv("DB_SSL_MODE", "disable"),
        },
        Auth: AuthConfig{
            JWTSecret:          mustGetEnv("JWT_SECRET"),
            JWTExpiry:          getDurationEnv("JWT_EXPIRY", 15*time.Minute),
            RefreshTokenExpiry: getDurationEnv("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
            BcryptCost:         getIntEnv("BCRYPT_COST", 12),
        },
        Storage: StorageConfig{
            UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
            MaxSize:    getInt64Env("MAX_UPLOAD_SIZE", 10*1024*1024), // 10MB
        },
    }

    return cfg, nil
}

func (c DatabaseConfig) DSN() string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// Helper functions
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func mustGetEnv(key string) string {
    value := os.Getenv(key)
    if value == "" {
        log.Fatalf("Environment variable %s is required", key)
    }
    return value
}

func getIntEnv(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
            return intVal
        }
    }
    return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
```

## 2. Simple Middleware Implementation

### internal/middleware/middleware.go
```go
package middleware

import (
    "context"
    "log"
    "net/http"
    "runtime/debug"
    "time"

    "github.com/google/uuid"
)

// Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain applies middlewares in order
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
    for i := len(middlewares) - 1; i >= 0; i-- {
        f = middlewares[i](f)
    }
    return f
}

// RequestID middleware
func RequestID(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        w.Header().Set("X-Request-ID", requestID)
        
        next(w, r.WithContext(ctx))
    }
}

// Logger middleware
func Logger(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer
        wrapped := &responseWriter{
            ResponseWriter: w,
            statusCode:    http.StatusOK,
        }
        
        next(wrapped, r)
        
        log.Printf(
            "%s %s %s %d %s %s",
            r.RemoteAddr,
            r.Method,
            r.RequestURI,
            wrapped.statusCode,
            time.Since(start),
            r.UserAgent(),
        )
    }
}

// Recovery middleware
func Recovery(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next(w, r)
    }
}

// CORS middleware
func CORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        
        next(w, r)
    }
}

// responseWriter wrapper
type responseWriter struct {
    http.ResponseWriter
    statusCode int
    written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
    if !rw.written {
        rw.statusCode = code
        rw.ResponseWriter.WriteHeader(code)
        rw.written = true
    }
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    if !rw.written {
        rw.WriteHeader(http.StatusOK)
    }
    return rw.ResponseWriter.Write(b)
}
```

### internal/middleware/auth.go
```go
package middleware

import (
    "context"
    "fmt"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

// Auth middleware
func Auth(jwtSecret string) Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            tokenString := extractToken(r)
            if tokenString == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("unexpected signing method")
                }
                return []byte(jwtSecret), nil
            })

            if err != nil || !token.Valid {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // Add user info to context
            ctx := r.Context()
            ctx = context.WithValue(ctx, "user_id", claims["user_id"])
            ctx = context.WithValue(ctx, "email", claims["email"])

            next(w, r.WithContext(ctx))
        }
    }
}

func extractToken(r *http.Request) string {
    bearerToken := r.Header.Get("Authorization")
    if strings.HasPrefix(bearerToken, "Bearer ") {
        return strings.TrimPrefix(bearerToken, "Bearer ")
    }
    return ""
}

// GetUserID from context
func GetUserID(ctx context.Context) string {
    if userID, ok := ctx.Value("user_id").(string); ok {
        return userID
    }
    return ""
}
```

## 3. Simple Server Setup

### internal/server/server.go
```go
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

    _ "github.com/lib/pq"
    
    "github.com/yourusername/lms/internal/config"
    "github.com/yourusername/lms/internal/database"
    "github.com/yourusername/lms/internal/handler"
    "github.com/yourusername/lms/internal/middleware"
    "github.com/yourusername/lms/internal/service"
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

func (s *Server) setupRoutes() http.Handler {
    mux := http.NewServeMux()

    // Initialize services
    authService := service.NewAuthService(s.queries, s.config)
    userService := service.NewUserService(s.queries)
    courseService := service.NewCourseService(s.queries)

    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService)
    userHandler := handler.NewUserHandler(userService)
    courseHandler := handler.NewCourseHandler(courseService)

    // Global middleware
    globalMiddleware := []middleware.Middleware{
        middleware.RequestID,
        middleware.Logger,
        middleware.Recovery,
        middleware.CORS,
    }

    // Auth middleware
    auth := middleware.Auth(s.config.Auth.JWTSecret)

    // Routes
    // Health check
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // Public routes
    mux.HandleFunc("/api/v1/auth/register", 
        chain(authHandler.Register, globalMiddleware...))
    mux.HandleFunc("/api/v1/auth/login", 
        chain(authHandler.Login, globalMiddleware...))
    mux.HandleFunc("/api/v1/auth/refresh", 
        chain(authHandler.RefreshToken, globalMiddleware...))

    // Protected routes
    mux.HandleFunc("/api/v1/users/me", 
        chain(userHandler.GetProfile, append(globalMiddleware, auth)...))
    mux.HandleFunc("/api/v1/users/me/update", 
        chain(userHandler.UpdateProfile, append(globalMiddleware, auth)...))

    // Course routes
    mux.HandleFunc("/api/v1/courses", courseHandler.HandleCourses)
    mux.HandleFunc("/api/v1/courses/", courseHandler.HandleCourse)

    return mux
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

// Helper function to chain middleware
func chain(f http.HandlerFunc, middlewares ...middleware.Middleware) http.HandlerFunc {
    return middleware.Chain(f, middlewares...)
}
```

## 4. Simple Handler Example

### internal/handler/auth.go
```go
package handler

import (
    "encoding/json"
    "net/http"

    "github.com/yourusername/lms/internal/service"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
    User         User   `json:"user"`
}

type User struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Call service
    result, err := h.authService.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // Send response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(LoginResponse{
        AccessToken:  result.AccessToken,
        RefreshToken: result.RefreshToken,
        User: User{
            ID:        result.User.ID,
            Email:     result.User.Email,
            FirstName: result.User.FirstName,
            LastName:  result.User.LastName,
        },
    })
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        Email     string `json:"email"`
        Password  string `json:"password"`
        FirstName string `json:"firstName"`
        LastName  string `json:"lastName"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate input
    if req.Email == "" || req.Password == "" {
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }

    // Call service
    user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.FirstName, req.LastName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Send response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "User registered successfully",
        "userId":  user.ID,
    })
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        RefreshToken string `json:"refreshToken"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Call service
    result, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // Send response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "accessToken":  result.AccessToken,
        "refreshToken": result.RefreshToken,
    })
}
```

### internal/handler/course.go
```go
package handler

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/yourusername/lms/internal/middleware"
    "github.com/yourusername/lms/internal/service"
)

type CourseHandler struct {
    courseService *service.CourseService
}

func NewCourseHandler(courseService *service.CourseService) *CourseHandler {
    return &CourseHandler{
        courseService: courseService,
    }
}

// HandleCourses handles /api/v1/courses
func (h *CourseHandler) HandleCourses(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.listCourses(w, r)
    case http.MethodPost:
        h.createCourse(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

// HandleCourse handles /api/v1/courses/{id}
func (h *CourseHandler) HandleCourse(w http.ResponseWriter, r *http.Request) {
    // Extract ID from path
    path := strings.TrimPrefix(r.URL.Path, "/api/v1/courses/")
    if path == "" {
        http.Error(w, "Course ID required", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case http.MethodGet:
        h.getCourse(w, r, path)
    case http.MethodPut:
        h.updateCourse(w, r, path)
    case http.MethodDelete:
        h.deleteCourse(w, r, path)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *CourseHandler) listCourses(w http.ResponseWriter, r *http.Request) {
    // Get query parameters
    query := r.URL.Query()
    page := query.Get("page")
    limit := query.Get("limit")
    search := query.Get("search")

    courses, err := h.courseService.ListCourses(r.Context(), page, limit, search)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(courses)
}

func (h *CourseHandler) createCourse(w http.ResponseWriter, r *http.Request) {
    // Check authentication
    userID := middleware.GetUserID(r.Context())
    if userID == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req struct {
        Title       string `json:"title"`
        Description string `json:"description"`
        Category    string `json:"category"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    course, err := h.courseService.CreateCourse(r.Context(), userID, req.Title, req.Description, req.Category)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) getCourse(w http.ResponseWriter, r *http.Request, id string) {
    course, err := h.courseService.GetCourse(r.Context(), id)
    if err != nil {
        http.Error(w, "Course not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) updateCourse(w http.ResponseWriter, r *http.Request, id string) {
    // Implementation
}

func (h *CourseHandler) deleteCourse(w http.ResponseWriter, r *http.Request, id string) {
    // Implementation
}
```

## 5. Main Entry Point

### cmd/api/main.go
```go
package main

import (
    "log"

    "github.com/yourusername/lms/internal/config"
    "github.com/yourusername/lms/internal/server"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Create server
    srv, err := server.New(cfg)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }

    // Start server
    if err := srv.Start(); err != nil {
        log.Fatal("Server failed:", err)
    }
}
```

## 6. Simple .env File

### .env
```env
# Server
PORT=8080
ENVIRONMENT=development
READ_TIMEOUT=15s
WRITE_TIMEOUT=15s

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=lms_user
DB_PASSWORD=your_password
DB_NAME=lms_db
DB_SSL_MODE=disable

# Auth
JWT_SECRET=your-super-secret-key-at-least-32-chars
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
BCRYPT_COST=12

# Storage
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=10485760
```

## 7. Makefile

```makefile
# Simple Makefile
.PHONY: run build test clean

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

migrate-up:
	goose -dir db/migrations postgres "postgres://lms_user:your_password@localhost/lms_db?sslmode=disable" up

migrate-down:
	goose -dir db/migrations postgres "postgres://lms_user:your_password@localhost/lms_db?sslmode=disable" down

sqlc:
	sqlc generate

dev: sqlc
	air

setup:
	go mod download
	go install github.com/cosmtrek/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Key Features of This Simple Implementation:

### 1. **Simple Configuration**
- Uses `godotenv` to load `.env` file
- Falls back to environment variables if `.env` not found
- Simple helper functions for different types

### 2. **Standard Library HTTP**
- Uses `http.ServeMux` for routing
- No external router dependencies
- Simple path-based routing

### 3. **Clean Middleware Pattern**
- Simple middleware chaining
- Easy to understand and extend
- Context-based value passing

### 4. **Practical Features**
- Request logging
- Panic recovery
- CORS support
- JWT authentication
- Request IDs

### 5. **Simple Handler Pattern**
- Method checking in handlers
- JSON request/response handling
- Error responses

## Usage:

```bash
# Create .env file
cp .env.example .env

# Install dependencies
go mod init github.com/yourusername/lms
go get github.com/joho/godotenv
go get github.com/golang-jwt/jwt/v5
go get github.com/google/uuid
go get github.com/lib/pq

# Run migrations
make migrate-up

# Generate SQLC
make sqlc

# Run the server
make run

# Or with hot reload
make dev
```

This implementation is:
- **Simple**: Easy to understand and modify
- **Standard**: Uses Go standard library for HTTP
- **Clean**: Clear separation of concerns
- **Practical**: Includes essential features for a web API
- **Scalable**: Can grow with your application needs