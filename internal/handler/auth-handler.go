package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net" // Add this
	"net/http"
	"time"

	"github.com/Abdelrahiim/lms/internal/config"
	"github.com/Abdelrahiim/lms/internal/database"
	"github.com/Abdelrahiim/lms/internal/middleware"
	"github.com/Abdelrahiim/lms/internal/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type AuthHandler struct {
	db      *sql.DB
	queries *database.Queries
	config  *config.Config
}

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

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func NewAuthHandler(db *sql.DB, queries *database.Queries, config *config.Config) *AuthHandler {
	return &AuthHandler{
		db:      db,
		queries: queries,
		config:  config,
	}
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

// Register handler with validation
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

	// Your registration logic here
	// The req variable contains the validated RegisterRequest
	// Create user
	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.SendErrorResponse(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

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
	// Return the created user data instead of just a message
	if err != nil {
		utils.SendErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(utils.SendMutationResponse("User created successfully")); err != nil {
		log.Fatal("failed to encode response: %w", err)
	}
}

// Login handler with validation
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

	// Check password
	err = utils.CheckPasswordHash(user.PasswordHash, req.Password)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, h.config.Auth.JWTSecret)
	if err != nil {
		utils.SendErrorResponse(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.SendErrorResponse(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	// Create session
	_, err = h.queries.CreateSession(r.Context(), database.CreateSessionParams{
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
		ExpiresAt:        time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		utils.SendErrorResponse(w, "Error creating session", http.StatusInternalServerError)
		return
	}

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

	jsonData, err := json.Marshal(loginResponse)
	if err != nil {
		utils.SendErrorResponse(w, "Error marshalling login response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jsonData); err != nil {
		log.Fatal("failed to encode response: %w", err)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetBearerToken(r.Header)
	if err != nil {
		utils.SendErrorResponse(w, "Error getting bearer token", http.StatusInternalServerError)
		return
	}
	claims, err := utils.ValidateJWT(token, h.config.Auth.JWTSecret)
	if err != nil {
		utils.SendErrorResponse(w, "Error validating token", http.StatusInternalServerError)
		return
	}

	session, err := h.queries.GetSession(r.Context(), database.GetSessionParams{
		UserID:     uuid.MustParse(claims.UserID),
		IpAddress:  pqtype.Inet{IPNet: net.IPNet{IP: net.ParseIP(utils.GetClientIP(r))}, Valid: true},
	})
	if err != nil {
		utils.SendErrorResponse(w, "Error getting session", http.StatusInternalServerError)
		return
	}

	err = h.queries.RevokeSession(r.Context(), database.RevokeSessionParams{
		ID:         session.ID,
		RevokedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		RevokedReason: sql.NullString{String: "User logged out", Valid: true},
	})
	if err != nil {
		utils.SendErrorResponse(w, "Error revoking session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(utils.SendMutationResponse("Logged out successfully")); err != nil {
		log.Fatal("failed to encode response: %w", err)
	}
}