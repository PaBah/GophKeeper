package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims - present JWT claims (customised with UserID)
type Claims struct {
	jwt.RegisteredClaims
	UserID    string
	SessionID string
}

// Parameters for JWT tokens generation/parsing
const (
	// TokenExp - JWT token expiration time in microseconds
	TokenExp = time.Hour * 3
)

// BuildJWTString - generate JWT string from UserID
func BuildJWTString(userID, sessionID, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID:    userID,
		SessionID: sessionID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenString, err
}

// GetUserID - parse JWT string and return UserID
func GetUserID(tokenString, secretKey string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}

	return claims.UserID
}

// GetSessionID - parse JWT string and return SessionID
func GetSessionID(tokenString, secretKey string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}

	return claims.SessionID
}

func IsValidToken(tokenString string, secret string) (bool, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	if !token.Valid {
		return false, fmt.Errorf("%w", err)
	}

	return true, nil
}
