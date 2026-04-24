package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error Hash password: %w", err)
	}
	return string(hashed), nil
}

func CheckHashedPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
	if err != nil {
		return fmt.Errorf("Password Do Not Match: %w", err)	
	}
	return nil	
}