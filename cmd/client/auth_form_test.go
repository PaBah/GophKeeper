package main

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Valid Email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Email without domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "Email without user",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "Empty Email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "Email with multiple '@'",
			email:   "test@example@.com",
			wantErr: true,
		},
		{
			name:    "Email with invalid characters",
			email:   "test@@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
