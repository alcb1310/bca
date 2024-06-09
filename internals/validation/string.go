package validation

import (
	"fmt"

	"github.com/badoux/checkmail"
)

const PASSWORD_SIZE = 3

func CompanyNameValidation(name string, required bool) string {
	if !required && name == "" {
		return ""
	}

	if required && name == "" {
		return "El nombre es requerido"
	}

	if len(name) < 3 {
		return "El nombre debe tener al menos 3 caracteres"
	}

	return ""
}

func EmailValidator(value string, required bool) string {
	if !required && value == "" {
		return ""
	}

	if required && value == "" {
		return "El email es requerido"
	}

	if err := checkmail.ValidateFormat(value); err != nil {
		return "Ingrese un correo válido"
	}

	return ""
}

func PasswordValidator(value string, required bool) string {
	if !required && value == "" {
		return ""
	}

	if required && value == "" {
		return "La contraseña es requerida"
	}

	if len(value) < PASSWORD_SIZE {
		return fmt.Sprintf("La contraseña debe tener al menos %d caracteres", PASSWORD_SIZE)
	}

	return ""
}

func NameValidation(value string, required bool) string {
	if !required && value == "" {
		return ""
	}

	if required && value == "" {
		return "es requerido"
	}

	return ""
}
