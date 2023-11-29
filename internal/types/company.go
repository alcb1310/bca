package types

import "github.com/google/uuid"

type Company struct {
	ID        uuid.UUID `json:"id"`
	Ruc       string    `json:"ruc"`
	Name      string    `json:"name"`
	Employees uint8     `json:"employees"`
	IsActive  bool      `json:"is_active"`
}
