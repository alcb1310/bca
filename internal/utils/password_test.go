package utils

import (
	"strings"
	"testing"
)

func TestComparePassword(t *testing.T) {
	t.Run("valid password", func(t *testing.T) {
		password := "password"
		hashed, _ := EncryptPasssword(password)

		valid, err := ComparePassword(string(hashed), password)

		if !valid {
			t.Errorf("Password '%s' should be valid", password)
		}

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		password := "password"
		hashed, _ := EncryptPasssword(password)

		valid, err := ComparePassword(string(hashed), "invalid")

		if valid {
			t.Errorf("Password '%s' should not be valid", password)
		}

		want := "Credenciales inválidas"
		got := err.Error()
		if strings.Compare(want, got) != 0 {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}
