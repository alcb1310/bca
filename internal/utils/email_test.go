package utils

import "testing"

func TestEmail(t *testing.T) {
	t.Run("Valid email", func(t *testing.T) {
		email := "wvXkz@example.com"

		if !IsValidEmail(email) {
			t.Errorf("Email '%s' should be valid", email)
		}
	})

	t.Run("Invalid email", func(t *testing.T) {
		email := "invalid"

		if IsValidEmail(email) {
			t.Errorf("Email '%s' should not be valid", email)
		}
	})

	t.Run("Empty email", func(t *testing.T) {
		email := ""

		if IsValidEmail(email) {
			t.Errorf("Email '%s' should not be valid", email)
		}
	})
}
