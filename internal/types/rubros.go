package types

import "github.com/google/uuid"

type Rubro struct {
	Id        uuid.UUID
	Code      string
	Name      string
	Unit      string
	CompanyId uuid.UUID
}
