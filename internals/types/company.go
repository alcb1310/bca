package types

import "github.com/google/uuid"

type Company struct {
	ID        uuid.UUID
	Ruc       string
	Name      string
	Employees uint
	IsActive  bool
}
