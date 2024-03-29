package types

import "github.com/google/uuid"

type ItemMaterialType struct {
	ItemId     uuid.UUID
	MaterialId uuid.UUID
	Quantity   float64
	CompanyId  uuid.UUID
}
