package types

import "github.com/google/uuid"

type BudgetItem struct {
	ID         uuid.UUID  `json:"id"`
	Code       string     `json:"code"`
	Name       string     `json:"name"`
	Level      uint8      `json:"level"`
	Accumulate *bool      `json:"accumulate"`
	ParentId   *uuid.UUID `json:"parent_id"`
	CompanyId  uuid.UUID  `json:"company_id"`
}
