package types

import (
	"database/sql"

	"github.com/google/uuid"
)

type BudgetItem struct {
	ID         uuid.UUID    `json:"id"`
	Code       string       `json:"code"`
	Name       string       `json:"name"`
	Level      uint8        `json:"level"`
	Accumulate sql.NullBool `json:"accumulate"`
	ParentId   *uuid.UUID   `json:"parent_id"`
	CompanyId  uuid.UUID    `json:"company_id"`
}

type BudgetItemResponse struct {
	ID         uuid.UUID      `json:"id"`
	Code       string         `json:"code"`
	Name       string         `json:"name"`
	Level      uint8          `json:"level"`
	Accumulate sql.NullBool   `json:"accumulate"`
	ParentId   uuid.NullUUID  `json:"parent_id"`
	ParentCode sql.NullString `json:"parent_code"`
	ParentName sql.NullString `json:"parent_name"`
	CompanyId  uuid.UUID      `json:"company_id"`
}

type BudgetItemJsonResponse struct {
	ID         uuid.UUID   `json:"id"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	Level      uint8       `json:"level"`
	Accumulate bool        `json:"accumulate"`
	Parent     *BudgetItem `json:"parent"`
	CompanyId  uuid.UUID   `json:"company_id"`
}
