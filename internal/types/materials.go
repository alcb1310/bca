package types

import "github.com/google/uuid"

type Material struct {
	Id        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Unit      string    `json:"unit"`
	Category  Category  `json:"category"`
	CompanyId uuid.UUID `json:"company_id"`
}
