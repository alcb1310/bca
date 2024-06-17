package validation

import "strconv"

func ValidatePositiveInteger(value string, required bool) (uint, string) {
	if !required && value == "" {
		return 0, ""
	}

	if required && value == "" {
		return 0, "es requerido"
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return 0, "debe ser un número entero"
	}

	if val < 0 {
		return 0, "debe ser un número positivo"
	}

	return uint(val), ""
}
