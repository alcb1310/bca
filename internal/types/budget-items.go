package types

import (
	"database/sql"

	"github.com/google/uuid"
)

type BudgetItem struct {
	ID         uuid.UUID
	Code       string
	Name       string
	Level      uint8
	Accumulate sql.NullBool
	ParentId   *uuid.UUID
	CompanyId  uuid.UUID
}

type BudgetItemResponse struct {
	ID         uuid.UUID
	Code       string
	Name       string
	Level      uint8
	Accumulate sql.NullBool
	ParentId   uuid.NullUUID
	ParentCode sql.NullString
	ParentName sql.NullString
	CompanyId  uuid.UUID
}
