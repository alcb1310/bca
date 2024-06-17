package validation

func ValidateLogin(fields map[string]string) map[string]string {
	validationErrors := make(map[string]string)

	if e := EmailValidator(fields["email"], true); e != "" {
		validationErrors["email"] = e
	}

	if e := PasswordValidator(fields["password"], true); e != "" {
		validationErrors["password"] = e
	}

	return validationErrors
}
