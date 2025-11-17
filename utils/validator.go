package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// ValidateEmail checks if email is valid
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	
	return nil
}

// ValidateUsername checks if username is valid
func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}
	
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	
	if len(username) > 30 {
		return errors.New("username must be less than 30 characters")
	}
	
	// Username should only contain alphanumeric characters, underscores, and hyphens
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, underscores, and hyphens")
	}
	
	return nil
}

// ValidatePassword checks if password meets security requirements
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
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
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}
	
	return nil
}

// ValidatePhone checks if phone number is valid (basic validation)
func ValidatePhone(phone string) error {
	if phone == "" {
		return nil // Phone is optional
	}
	
	// Remove spaces and dashes
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	
	// Basic phone validation (should start with + and contain 7-15 digits)
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{6,14}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone number format")
	}
	
	return nil
}

// SanitizeString removes potentially harmful characters
func SanitizeString(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	return input
}

