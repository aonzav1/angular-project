package user

import (
	"golang.org/x/crypto/bcrypt"
)

// Hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	// Generate a salt with default cost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Compares a plain text password with a hashed password
func ComparePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
