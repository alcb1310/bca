package types

import "github.com/google/uuid"

type BudgetItem struct {
	ID         uuid.UUID
	Code       string
	Name       string
	Level      uint8
	Accumulate *bool
	ParentId   *uuid.UUID
	CompanyId  uuid.UUID
}

type BudgetItemResponse struct {
	ID         uuid.UUID
	Code       string
	Name       string
	Level      uint8
	Accumulate *bool
	ParentId   *uuid.UUID
	ParentCode *string
	ParentName *string
	CompanyId  uuid.UUID
}
