package validation

import (
	"fmt"

	"github.com/alcb1310/bca/internals/types"
	"github.com/alcb1310/bca/internals/utils"
)

func ValidateRegisterForm(fields map[string]string, company *types.Company, user *types.CreateUser) map[string]string {
	validationErrors := make(map[string]string)

	if e := IdValidation(fields["ruc"], true); e != "" {
		validationErrors["ruc"] = e
	}
	company.Ruc = fields["ruc"]

	if e := CompanyNameValidation(fields["name"], true); e != "" {
		validationErrors["name"] = e
	}
	company.Name = fields["name"]

	emp, e := ValidatePositiveInteger(fields["employees"], false)
	if e != "" {
		validationErrors["employees"] = fmt.Sprintf("Empleados debe ser %s", e)
	}
	if emp == 0 {
		emp = 1
	}
	company.Employees = emp

	if e := EmailValidator(fields["email"], true); e != "" {
		validationErrors["email"] = e
	}
	user.Email = fields["email"]

	if e := PasswordValidator(fields["password"], true); e != "" {
		validationErrors["password"] = e
	}
	user.Password = utils.EncryptPassword(fields["password"])

	if e := NameValidation(fields["username"], true); e != "" {
		validationErrors["username"] = fmt.Sprintf("El nombre %s", e)
	}
	user.Name = fields["username"]

	return validationErrors
}
