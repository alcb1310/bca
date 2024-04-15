package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateUUID(t *testing.T) {
	t.Run("valid uuid", func(t *testing.T) {
		want := uuid.New()
		strUUUID := want.String()

		got, err := ValidateUUID(strUUUID, "proyecto")

		if err != nil {
			t.Errorf("ValidateUUID() error = %v", err)
			return
		}

		if got != want {
			t.Errorf("ValidateUUID() = %v, want %v", got, want)
		}
	})

	t.Run("invalid uuid", func(t *testing.T) {
		strUUID := "f3b8b49a-1a1d-4252-be67-z28d54c87eae"
		_, err := ValidateUUID(strUUID, "proyecto")
		if err == nil {
			t.Errorf("Expected an error and got none")
			return
		}

		want := "invalid UUID format"
		if err.Error() != want {
			t.Errorf("Expected error to be '%s' but got '%s'", want, err.Error())
		}
	})

	t.Run("invalid uuid length", func(t *testing.T) {
		strUUID := "1a1d-4252-be67-238d54c87eae"
		_, err := ValidateUUID(strUUID, "proyecto")
		if err == nil {
			t.Errorf("Expected an error and got none")
			return
		}

		want := "invalid UUID length: 27"
		if err.Error() != want {
			t.Errorf("Expected error to be '%s' but got '%s'", want, err.Error())
		}
	})

	t.Run("empty uuid", func(t *testing.T) {
		strUUID := ""
		_, err := ValidateUUID(strUUID, "proyecto")
		if err == nil {
			t.Errorf("Expected an error and got none")
			return
		}

		want := "Seleccione un proyecto"
		if err.Error() != want {
			t.Errorf("Expected error to be '%s' but got '%s'", want, err.Error())
		}
	})
}
