package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, cost int) (string, error) {
	if password == "" {
		return "", errors.New("password is empty")
	}
	if cost == 0 {
		cost = bcrypt.DefaultCost 
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
