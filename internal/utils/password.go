package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengenkripsi password menggunakan bcrypt
func HashPassword(password string) (string, error) {
	// Cost 10 adalah default dan cukup secure
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword membandingkan password dengan hash
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
