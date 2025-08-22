package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Abdelrahiim/lms/internal/config"
	"github.com/Abdelrahiim/lms/internal/database"
	"github.com/Abdelrahiim/lms/internal/middleware"
	"github.com/Abdelrahiim/lms/internal/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// ============================================================================
// TYPES AND STRUCTS
// ============================================================================

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	db      *sql.DB
	queries *database.Queries
	config  *config.Config
}

// RegisterRequest represents the user registration payload
type RegisterRequest struct {
	Email             string `json:"email" validate:"required,email"`
	Password          string `json:"password" validate:"required,min=8"`
	FirstName         string `json:"firstName" validate:"required"`
	LastName          string `json:"lastName" validate:"required"`
	DisplayName       string `json:"displayName,omitempty"`
	AvatarURL         string `json:"avatarUrl,omitempty" validate:"omitempty,url"`
	Bio               string `json:"bio,omitempty" validate:"omitempty,max=500"`
	Phone             string `json:"phone,omitempty" validate:"omitempty,e164"`
	DateOfBirth       string `json:"dateOfBirth,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Gender            string `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Country           string `json:"country,omitempty" validate:"omitempty,len=2"`
	Timezone          string `json:"timezone,omitempty" validate:"omitempty,timezone"`
	PreferredLanguage string `json:"preferredLanguage,omitempty" validate:"omitempty,len=2"`
}

// LoginRequest represents the user login payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the successful login response
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

// RefreshResponse represents the token refresh response
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

// User represents the user data in responses
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(db *sql.DB, queries *database.Queries, config *config.Config) *AuthHandler {
	return &AuthHandler{
		db:      db,
		queries: queries,
		config:  config,
	}
}

// ============================================================================
// HTTP HANDLERS
// ============================================================================

// Register handles user registration requests
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Get validated payload from context
	req, ok := middleware.GetValidatedPayload[RegisterRequest](r)
	if !ok {
		return
	}

	// Check if email already exists
	_, err := h.queries.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		// Email already exists
		utils.SendErrorResponse(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.SendErrorResponse(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Create user with all provided data
	err = h.queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		DisplayName:  sql.NullString{String: req.DisplayName, Valid: req.DisplayName != ""},
		AvatarUrl:    sql.NullString{String: req.AvatarURL, Valid: req.AvatarURL != ""},
		Bio:          sql.NullString{String: req.Bio, Valid: req.Bio != ""},
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		DateOfBirth: func() sql.NullTime {
			if req.DateOfBirth == "" {
				return sql.NullTime{Valid: false}
			}
			t, err := time.Parse("2006-01-02", req.DateOfBirth)
			if err != nil {
				return sql.NullTime{Valid: false}
			}
			return sql.NullTime{Time: t, Valid: true}
		}(),
		Gender:            sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Country:           sql.NullString{String: req.Country, Valid: req.Country != ""},
		Timezone:          sql.NullString{String: req.Timezone, Valid: req.Timezone != ""},
		PreferredLanguage: sql.NullString{String: req.PreferredLanguage, Valid: req.PreferredLanguage != ""},
	})
	if err != nil {
		utils.SendErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(utils.SendMutationResponse("User created successfully")); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// Login handles user login requests
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Get validated payload from context
	req, ok := middleware.GetValidatedPayload[LoginRequest](r)
	if !ok {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if user exists
	user, err := h.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	err = utils.CheckPasswordHash(user.PasswordHash, req.Password)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, h.config.Auth.JWTSecret)
	if err != nil {
		utils.SendErrorResponse(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.SendErrorResponse(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	// Create user session with device and location information
	err = h.queries.CreateSession(r.Context(), database.CreateSessionParams{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: refreshToken,
		AccessTokenHash:  sql.NullString{String: accessToken, Valid: true},
		DeviceName:       sql.NullString{String: r.Header.Get("User-Agent"), Valid: r.Header.Get("User-Agent") != ""},
		DeviceType:       sql.NullString{String: utils.GetDeviceType(r), Valid: true},
		Browser:          sql.NullString{String: utils.GetBrowser(r), Valid: true},
		BrowserVersion:   sql.NullString{String: utils.GetBrowserVersion(r), Valid: true},
		Os:               sql.NullString{String: utils.GetOS(r), Valid: true},
		OsVersion:        sql.NullString{String: utils.GetOSVersion(r), Valid: true},
		IpAddress:        pqtype.Inet{IPNet: net.IPNet{IP: net.ParseIP(utils.GetClientIP(r))}, Valid: true},
		Location:         pqtype.NullRawMessage{RawMessage: []byte(utils.GetLocation(r)), Valid: utils.GetLocation(r) != ""},
		IsActive:         sql.NullBool{Bool: true, Valid: true},
		LastAccessedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		ExpiresAt:        time.Now().AddDate(0, 0, 60), // 60 days expiration
	})
	if err != nil {
		utils.SendErrorResponse(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	// Prepare login response
	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: User{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	// Send login response
	jsonData, err := json.Marshal(loginResponse)
	if err != nil {
		utils.SendErrorResponse(w, "Error marshalling login response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jsonData); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// Logout handles user logout requests
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Extract and validate bearer token
	token, err := utils.GetBearerToken(r.Header)
	if err != nil {
		utils.SendErrorResponse(w, "Error getting bearer token", http.StatusUnauthorized)
		return
	}

	// Validate JWT token
	claims, err := utils.ValidateJWT(token, h.config.Auth.JWTSecret)
	if err != nil {
		utils.SendErrorResponse(w, "Error validating token", http.StatusUnauthorized)
		return
	}

	// Get user session by user ID and IP address
	session, err := h.queries.GetSessionByUserID(r.Context(), database.GetSessionByUserIDParams{
		UserID:    uuid.MustParse(claims.UserID),
		IpAddress: pqtype.Inet{IPNet: net.IPNet{IP: net.ParseIP(utils.GetClientIP(r))}, Valid: true},
	})
	if err != nil {
		utils.SendErrorResponse(w, "Error getting session", http.StatusInternalServerError)
		return
	}

	// Revoke the session
	err = h.queries.RevokeSession(r.Context(), database.RevokeSessionParams{
		ID:            session.ID,
		RevokedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		RevokedReason: sql.NullString{String: "User logged out", Valid: true},
	})
	if err != nil {
		utils.SendErrorResponse(w, "Error revoking session", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(utils.SendMutationResponse("Logged out successfully")); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// Refresh handles token refresh requests
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from request
	refreshToken, err := utils.GetBearerToken(r.Header)
	if err != nil {
		utils.SendErrorResponse(w, "Error getting bearer token", http.StatusUnauthorized)
		return
	}

	// Get session by refresh token
	session, err := h.queries.GetSessionByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Get user information
	user, err := h.queries.GetUser(r.Context(), session.UserID)
	if err != nil {
		utils.SendErrorResponse(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, h.config.Auth.JWTSecret)
	if err != nil {
		utils.SendErrorResponse(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	// Update session last accessed time
	err = h.queries.UpdateSessionLastAccessedAt(r.Context(), database.UpdateSessionLastAccessedAtParams{
		ID:             session.ID,
		LastAccessedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		utils.SendErrorResponse(w, "Error updating session last accessed at", http.StatusInternalServerError)
		return
	}

	// Prepare refresh response
	refreshResponse := RefreshResponse{
		AccessToken: accessToken,
	}

	// Send refresh response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(refreshResponse); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
