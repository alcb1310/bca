package types

import "github.com/google/uuid"

type Rubro struct {
	Id        uuid.UUID
	Code      string
	Name      string
	Unit      string
	CompanyId uuid.UUID
}

type ACU struct {
	Item      Rubro
	Material  Material
	Quantity  float64
	CompanyId uuid.UUID
}
