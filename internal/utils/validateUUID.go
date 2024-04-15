package utils

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ValidateUUID returns a valid uuid or an error
func ValidateUUID(id, name string) (uuid.UUID, error) {
	if id == "" {
		return uuid.UUID{}, errors.New(fmt.Sprintf("Seleccione un %s", name))
	}
	return uuid.Parse(id)
}
