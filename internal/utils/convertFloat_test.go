package utils

import "testing"

func TestConvertFloat(t *testing.T) {
	t.Run("Valid float", func(t *testing.T) {
		val := "10.25"
		want := 10.25

		got, err := ConvertFloat(val, "test", true)
		if err != nil {
			t.Errorf("ConvertFloat() error = %v", err)
			return
		}

		if got != want {
			t.Errorf("ConvertFloat() = %v, want %v", got, want)
		}
	})

	t.Run("Invalid float", func(t *testing.T) {
		val := "invalid"
		want := "test debe ser un número válido"

		_, err := ConvertFloat(val, "test", true)
		if err == nil {
			t.Errorf("Expected an error and got none")
			return
		}

		if err.Error() != want {
			t.Errorf("Expected error to be '%s' but got '%s'", want, err.Error())
		}

	})

	t.Run("Required empty value", func(t *testing.T) {
		val := ""
		want := "test es requerido"

		_, err := ConvertFloat(val, "test", true)
		if err == nil {
			t.Errorf("Expected an error and got none")
			return
		}

		if err.Error() != want {
			t.Errorf("Expected error to be '%s' but got '%s'", want, err.Error())
		}
	})

	t.Run("Not required not empty value", func(t *testing.T) {
		val := ""
		want := 0.0

		got, err := ConvertFloat(val, "test", false)
		if err != nil {
			t.Errorf("ConvertFloat() error = %v", err)
			return
		}

		if got != want {
			t.Errorf("ConvertFloat() = %v, want %v", got, want)
		}
	})

	t.Run("El valor debe ser postivo", func(t *testing.T) {
		val := "-1.25"
		want := "test debe ser un número positivo"

		_, err := ConvertFloat(val, "test", true)
		if err == nil {
			t.Errorf("Expected an error and got none")
			return
		}

		if err.Error() != want {
			t.Errorf("Expected error to be '%s' but got '%s'", want, err.Error())
		}
	})
}
