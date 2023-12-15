package types

import "github.com/google/uuid"

type Project struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  *bool     `json:"is_active"`
	CompanyId uuid.UUID `json:"company_id"`
}
