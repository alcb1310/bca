package types

import "github.com/google/uuid"

type ItemMaterialType struct {
	ItemId     uuid.UUID `json:"item_id"`
	MaterialId uuid.UUID `json:"material_id"`
	Quantity   float64   `json:"quantity"`
	CompanyId  uuid.UUID `json:"company_id"`
}
