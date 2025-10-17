package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims adalah struktur claims untuk JWT token
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT token baru
func GenerateToken(userID, username, role, secret string, expiration time.Duration) (string, error) {
	// Buat claims
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken memvalidasi JWT token dan return claims
func ValidateToken(tokenString, secret string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Pastikan signing method adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Ambil claims dari token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// ParseDuration helper untuk parse duration string
func ParseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return 24 * time.Hour // default 24 hours
	}
	return d
}
