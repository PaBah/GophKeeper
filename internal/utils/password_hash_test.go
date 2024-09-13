package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHash(t *testing.T) {
	tt := []struct {
		name      string
		password  string
		hashedErr bool
	}{
		{"EmptyString", "", false},
		{"NormalPassword", "password123", false},
		{"SpecialChars", "password!@#", false},
		{"UnicodeChars", "passwordΣ", false},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword := PasswordHash(tc.password)
			err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tc.password))
			if err != nil && !tc.hashedErr {
				t.Fatalf("Password and hashed password do not match. Error: %v", err)
			}
			if err == nil && tc.hashedErr {
				t.Fatalf("Expected error but received nil while creating Password Hash.")
			}
		})
	}
}
func TestCheckPasswordHash(t *testing.T) {
	tt := []struct {
		name           string
		password       string
		hashedPassword string
		expectedResult bool
	}{
		{"Matching", "password123", PasswordHash("password123"), true},
		{"Non-Matching", "password123", PasswordHash("password1234"), false},
		{"SpecialChars", "password!@#", PasswordHash("password!@#"), true},
		{"UnicodeChars", "passwordΣ", PasswordHash("passwordΣ"), true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			isMatch := CheckPasswordHash(tc.hashedPassword, tc.password)
			if isMatch != tc.expectedResult {
				t.Fatalf("Expected %v but received %v for Password '%s' and Hashed Password '%s'", tc.expectedResult, isMatch, tc.password, tc.hashedPassword)
			}
		})
	}
}
