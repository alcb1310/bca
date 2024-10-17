package types

import "github.com/google/uuid"

type Rubro struct {
	Id        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Unit      string    `json:"unit"`
	CompanyId uuid.UUID `json:"company_id"`
}

type ACU struct {
	Item      Rubro     `json:"item"`
	Material  Material  `json:"material"`
	Quantity  float64   `json:"quantity"`
	CompanyId uuid.UUID `json:"company_id"`
}
