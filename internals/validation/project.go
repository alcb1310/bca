package validation

import "github.com/alcb1310/bca/internals/types"

func ValidateProject(fields map[string]string) (types.Project, map[string]string) {
	project := types.Project{}
	validationErrors := make(map[string]string)

	if e := NameValidation(fields["name"], true); e != "" {
		validationErrors["name"] = e
	}
	project.Name = fields["name"]

	return project, validationErrors
}
