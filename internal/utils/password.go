package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPasssword(p string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(p), 8)
}

func ComparePassword(hashed, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return false, errors.New("Credenciales inválidas")
	}

	return true, nil
}
