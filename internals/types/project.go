package types

import "github.com/google/uuid"

type Project struct {
	ID        uuid.UUID
	Name      string
	IsActive  bool
	CompanyID uuid.UUID
}
