package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	apperrors "backend-service-internpro/internal/pkg/errors"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// ValidationResult holds validation errors
type ValidationResult struct {
	Errors map[string]string `json:"errors,omitempty"`
}

// HasErrors returns true if there are validation errors
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// AddError adds a validation error
func (vr *ValidationResult) AddError(field, message string) {
	if vr.Errors == nil {
		vr.Errors = make(map[string]string)
	}
	vr.Errors[field] = message
}

// ToAppError converts validation result to AppError
func (vr *ValidationResult) ToAppError() *apperrors.AppError {
	var details []string
	for field, message := range vr.Errors {
		details = append(details, fmt.Sprintf("%s: %s", field, message))
	}
	return apperrors.ValidationFailed(strings.Join(details, ", "))
}

// Validator provides validation methods
type Validator struct{}

// New creates a new validator instance
func New() *Validator {
	return &Validator{}
}

// ValidateLoginRequest validates login request
func (v *Validator) ValidateLoginRequest(req interface{}) *ValidationResult {
	// Type assertion based on your DTO
	result := &ValidationResult{}

	// This would be implemented based on your specific DTO structure
	// For now, returning empty result
	return result
}

// Email validation
func (v *Validator) IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// Password validation
func (v *Validator) IsValidPassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	if len(password) > 128 {
		return false, "Password must be less than 128 characters"
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "Password must contain at least one number"
	}
	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	return true, ""
}

// Username validation
func (v *Validator) IsValidUsername(username string) (bool, string) {
	if len(username) < 3 {
		return false, "Username must be at least 3 characters long"
	}
	if len(username) > 60 {
		return false, "Username must be less than 60 characters"
	}

	// Allow alphanumeric and underscore
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != '_' {
			return false, "Username can only contain letters, numbers, and underscore"
		}
	}

	return true, ""
}

// OTP validation
func (v *Validator) IsValidOTP(otp string) (bool, string) {
	if len(otp) != 6 {
		return false, "OTP must be exactly 6 digits"
	}

	for _, char := range otp {
		if !unicode.IsNumber(char) {
			return false, "OTP must contain only numbers"
		}
	}

	return true, ""
}

// Generic string validation
func (v *Validator) IsRequired(value, fieldName string) (bool, string) {
	if strings.TrimSpace(value) == "" {
		return false, fmt.Sprintf("%s is required", fieldName)
	}
	return true, ""
}

func (v *Validator) IsMaxLength(value string, maxLen int, fieldName string) (bool, string) {
	if len(value) > maxLen {
		return false, fmt.Sprintf("%s must be less than %d characters", fieldName, maxLen)
	}
	return true, ""
}

func (v *Validator) IsMinLength(value string, minLen int, fieldName string) (bool, string) {
	if len(value) < minLen {
		return false, fmt.Sprintf("%s must be at least %d characters", fieldName, minLen)
	}
	return true, ""
}
