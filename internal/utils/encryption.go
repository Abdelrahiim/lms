package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CustomClaims represents the JWT claims structure following industry standards
type CustomClaims struct {
	UserID string `json:"sub"`           // Subject (user ID)
	Email  string `json:"email"`         // User email
	Role   string `json:"role,omitempty"` // User role (optional)
	Type   string `json:"typ"`           // Token type (access_token)
	jwt.RegisteredClaims
}

// HashPassword creates a bcrypt hash of the provided password
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent
func CheckPasswordHash(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match: %v", err)
	}
	return nil
}

// GenerateAccessToken creates a JWT access token with standard claims
func GenerateAccessToken(userID uuid.UUID, email string, secretKey string) (string, error) {
	now := time.Now().UTC()
	
	claims := CustomClaims{
		UserID: userID.String(),
		Email:  email,
		Type:   "access_token",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    "lms-api",
			Audience:  jwt.ClaimStrings{"lms-web", "lms-mobile"},
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error generating access token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT parses and validates a JWT token, returning the claims if valid
func ValidateJWT(tokenString string, tokenSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate user ID format
	if _, err := uuid.Parse(claims.UserID); err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	return claims, nil
}

// GetUserIDFromClaims extracts and parses the user ID from validated JWT claims
func GetUserIDFromClaims(claims *CustomClaims) (uuid.UUID, error) {
	return uuid.Parse(claims.UserID)
}

// GenerateRefreshToken creates a cryptographically secure random refresh token
func GenerateRefreshToken() (string, error) {
	// Generate 64 random bytes for stronger security
	tokenBytes := make([]byte, 64)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("error generating refresh token: %v", err)
	}

	// Encode to base64 URL-safe format
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	return token, nil
}

// GetBearerToken extracts the bearer token from the Authorization header
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no Authorization header found")
	}
	
	if len(authHeader) < 7 || authHeader[:6] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization header format")
	}
	
	return authHeader[7:], nil
}
