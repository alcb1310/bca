package utils

import (
	"net/mail"
)

func IsValidEmail(e string) bool {
	_, err := mail.ParseAddress(e)

	return err == nil
}
