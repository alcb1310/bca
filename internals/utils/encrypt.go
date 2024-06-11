package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const COST = 14

func EncryptPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), COST)
	return string(bytes)
}

func IsValidPassword(hashedPasssword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasssword), []byte(password))

	return err != nil
}
