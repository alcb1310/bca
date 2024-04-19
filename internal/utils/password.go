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

	return err == nil, errors.New("Credenciales inválidas")
}
