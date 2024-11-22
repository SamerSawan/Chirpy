package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", fmt.Errorf("Unable to hash password: %w", err)
	}
	return string(hashed_pass), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("Incorrect password")
	}
	return nil
}
