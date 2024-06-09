package validation

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
