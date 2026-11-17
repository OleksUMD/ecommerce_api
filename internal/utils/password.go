package utils

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword created encrypted hash of a raw password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a raw password with its encrypted version
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func ValidatePassword(password string) error {
	// Validate password length, required symbols etc.
	minLen := 8
	maxLen := 16
	isCapital := false
	isSpec := false
	specSymbs := "!@_?#*()[]$%"

	if len(password) < minLen {
		return fmt.Errorf("password is too short. Min length is %d characters", minLen)
	} else if len(password) > maxLen {
		return fmt.Errorf("password is too long. Max length is %d characters", maxLen)
	}
	for _, c := range password {
		if c > unicode.MaxASCII {
			return errors.New("only ascii characters accepted")
		}
		if unicode.IsSpace(c) {
			return errors.New("space is not accepted")
		}
		switch {
		case unicode.IsUpper(c):
			isCapital = true
		case strings.Contains(specSymbs, string(c)):
			isSpec = true
		}
	}
	if !isCapital {
		return errors.New("password should contain at least one upper case letter")
	}
	if !isSpec {
		return fmt.Errorf("password should contain at least one special character: '%s'", specSymbs)
	}
	return nil
}
