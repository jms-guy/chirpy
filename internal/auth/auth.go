package auth

import (
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

func HashPassword(password string) (string, error) {	//Takes a password string and hashes it
	pass := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return "Error hashing password", err
	}
	return string(hash), nil
}

func CheckPasswordHash(hash, password string) error {	//Checks a given password against a hash 
	pass := []byte(password)

	err := bcrypt.CompareHashAndPassword([]byte(hash), pass)
	if err != nil {
		return fmt.Errorf("error comparing password against hash: %w", err)
	}
	return nil
}