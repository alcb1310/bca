package types

import "github.com/google/uuid"

type Quantity struct {
	Id        uuid.UUID
	Project   Project
	Rubro     Rubro
	Quantity  float64
	CompanyId uuid.UUID
}
