package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(newPassword string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match: %v", err)
	}
	return nil
}

func GenerateAccessToken(userID uuid.UUID, email string, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
		"type":    "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error generating access token: %v", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {
	// Generate 32 random bytes for the refresh token
	// Generate a cryptographically secure random token using crypto/rand
	tokenBytes := make([]byte, 64) // Using 64 bytes for stronger security
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("error generating refresh token: %v", err)
	}

	// Encode the random bytes to base64 to make it URL-safe
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	// Append user ID to allow validation, using a more secure delimiter

	return token, nil
}
