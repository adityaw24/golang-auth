package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

//Encrypt plain password
func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("| Failed to hash a password "+err.Error())
		panic("Failed to hash a password")
	}
	return string(hash), err
}

// Compare plain hashed password retrieved from db against user-entered password
func ComparePassword(hashedPass string, plainPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("| Incorrect password. Error "+err.Error())
		return false
	}
	log.Println("| Password Matched.")
	return true
}