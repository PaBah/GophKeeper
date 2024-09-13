package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestBuildJWTString(t *testing.T) {
	secretKey := "mysecret"
	userID := "user1"
	sessionID := "session1"

	tt := []struct {
		name string
	}{
		{name: "Normal"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tokenStr, err := BuildJWTString(userID, sessionID, secretKey)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) { return []byte(secretKey), nil })

			if claims, ok := token.Claims.(*Claims); ok && token.Valid {
				if claims.UserID != userID {
					t.Errorf("UserID mismatch. Expected: %s, Got: %s", userID, claims.UserID)
				}

				if claims.SessionID != sessionID {
					t.Errorf("SessionID mismatch. Expected: %s, Got: %s", sessionID, claims.SessionID)
				}
			} else {
				t.Errorf("Failed to parse token: %v", err)
			}

			_, err = BuildJWTString(userID, sessionID, "")
			if err != nil {
				t.Errorf("Expected error due to empty secretKey, but got nil")
			}
		})
	}
}
func TestGetUserID(t *testing.T) {
	secretKey := "mysecret"
	userID := "user1"
	sessionID := "session1"

	token, _ := BuildJWTString(userID, sessionID, secretKey)
	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, -2)),
		},
		UserID:    userID,
		SessionID: sessionID,
	})

	invalidTokenString, _ := invalidToken.SignedString([]byte(secretKey))
	tt := []struct {
		name      string
		tokenStr  string
		secretKey string
		expected  string
	}{
		{name: "Normal", tokenStr: token, secretKey: "mysecret", expected: "user1"},
		{name: "Wrong Secret Key", tokenStr: token, secretKey: "wrongSecret", expected: ""},
		{name: "Expired token", tokenStr: invalidTokenString, secretKey: "mysecret", expected: ""},
		{name: "Invalid Token", tokenStr: "invalidToken", secretKey: "mySecret", expected: ""},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetUserID(tc.tokenStr, tc.secretKey)
			if actual != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, actual)
			}
		})
	}
}

func TestGetSessionID(t *testing.T) {
	secretKey := "mysecret"
	userID := "user1"
	sessionID := "session1"

	token, _ := BuildJWTString(userID, sessionID, secretKey)
	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, -2)),
		},
		UserID:    userID,
		SessionID: sessionID,
	})

	invalidTokenString, _ := invalidToken.SignedString([]byte(secretKey))
	tt := []struct {
		name      string
		tokenStr  string
		secretKey string
		expected  string
	}{
		{name: "Normal", tokenStr: token, secretKey: "mysecret", expected: "session1"},
		{name: "Wrong Secret Key", tokenStr: token, secretKey: "wrongSecret", expected: ""},
		{name: "Expired token", tokenStr: invalidTokenString, secretKey: "mysecret", expected: ""},
		{name: "Invalid Token", tokenStr: "invalidToken", secretKey: "mySecret", expected: ""},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetSessionID(tc.tokenStr, tc.secretKey)
			if actual != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, actual)
			}
		})
	}
}

func TestIsValidToken(t *testing.T) {
	secretKey := "mysecret"
	userID := "user1"
	sessionID := "session1"

	token, _ := BuildJWTString(userID, sessionID, secretKey)
	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, -2)),
		},
		UserID:    userID,
		SessionID: sessionID,
	})

	invalidTokenString, _ := invalidToken.SignedString([]byte(secretKey))
	tt := []struct {
		name      string
		tokenStr  string
		secretKey string
		expected  bool
	}{
		{name: "Valid Token", tokenStr: token, secretKey: "mysecret", expected: true},
		{name: "Expired Token", tokenStr: invalidTokenString, secretKey: "mysecret", expected: false},
		{name: "Invalid Token", tokenStr: "invalidToken", secretKey: "mySecret", expected: false},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, _ := IsValidToken(tc.tokenStr, tc.secretKey)
			if actual != tc.expected {
				t.Errorf("Expected: %v, Got: %v", tc.expected, actual)
			}
		})
	}
}
