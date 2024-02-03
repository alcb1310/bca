package types

import "github.com/google/uuid"

type Project struct {
	ID        uuid.UUID
	Name      string
	IsActive  *bool
	CompanyId uuid.UUID
	GrossArea float64
	NetArea   float64
}
