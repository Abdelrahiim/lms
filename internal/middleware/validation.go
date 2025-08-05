package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a validation error response
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// Validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function to use JSON field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	registerCustomValidators()
}

// registerCustomValidators adds custom validation rules
func registerCustomValidators() {
	// Custom timezone validator
	validate.RegisterValidation("timezone", func(fl validator.FieldLevel) bool {
		tz := fl.Field().String()
		if tz == "" {
			return true // Allow empty for optional fields
		}
		// Add your timezone validation logic here
		// For now, just check if it's not empty
		return len(tz) > 0
	})
}

// ValidateJSON middleware validates JSON request body against a struct
func ValidateJSON[T any](next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
			next(w, r)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
			return
		}

		var payload T
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Strict JSON parsing

		if err := decoder.Decode(&payload); err != nil {
			responseError := ErrorResponse{
				Error:   "invalid_json",
				Message: "Invalid JSON format",
			}
			respondWithError(w, http.StatusBadRequest, responseError)
			return
		}

		// Validate the struct
		if err := validate.Struct(payload); err != nil {
			validationErrors := []ValidationError{}

			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, ValidationError{
					Field:   err.Field(),
					Tag:     err.Tag(),
					Value:   fmt.Sprintf("%v", err.Value()),
					Message: getErrorMessage(err),
				})
			}

			responseError := ErrorResponse{
				Error:   "validation_failed",
				Message: "Request validation failed",
				Errors:  validationErrors,
			}
			respondWithError(w, http.StatusBadRequest, responseError)
			return
		}

		// Store validated payload in context for handler use
		type validatedPayloadKey struct{}
		ctx := r.Context()
		ctx = context.WithValue(ctx, validatedPayloadKey{}, payload)

		next(w, r.WithContext(ctx))
	}
}

// GetValidatedPayload retrieves the validated payload from context
func GetValidatedPayload[T any](r *http.Request) (T, bool) {
	type validatedPayloadKey struct{}
	value := r.Context().Value(validatedPayloadKey{})
	if value == nil {
		var zero T
		return zero, false
	}
	payload, ok := value.(T)
	return payload, ok
}

// getErrorMessage returns user-friendly error messages
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", err.Field(), err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", err.Field(), err.Param())
	case "datetime":
		return fmt.Sprintf("%s must be a valid date in format %s", err.Field(), err.Param())
	case "timezone":
		return fmt.Sprintf("%s must be a valid timezone", err.Field())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// respondWithError sends JSON error response
func respondWithError(w http.ResponseWriter, statusCode int, errorResponse ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
