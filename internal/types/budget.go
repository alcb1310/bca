package types

import (
	"database/sql"

	"github.com/google/uuid"
)

type Budget struct {
	ProjectId    uuid.UUID
	BudgetItemId uuid.UUID

	InitialQuantity sql.NullFloat64
	InitialCost     sql.NullFloat64
	InitialTotal    float64

	SpentQuantity sql.NullFloat64
	SpentTotal    float64

	RemainingQuantity sql.NullFloat64
	RemainingCost     sql.NullFloat64
	RemainingTotal    float64

	UpdatedBudget float64

	CompanyId uuid.UUID
}

// Structure for creating budget
type CreateBudget struct {
	ProjectId    uuid.UUID `json:"project_id"`
	BudgetItemId uuid.UUID `json:"budget_item_id"`
	Quantity     float64   `json:"quantity"`
	Cost         float64   `json:"cost"`
	CompanyId    uuid.UUID `json:"company_id"`
}

// Structure for reading budgets
type GetBudget struct {
	Project    ProjectData    `json:"project"`
	BudgetItem BudgetItemData `json:"budget_item"`

	InitialQuantity sql.NullFloat64 `json:"initial_quantity"`
	InitialCost     sql.NullFloat64 `json:"initial_cost"`
	InitialTotal    float64         `json:"initial_total"`

	SpentQuantity sql.NullFloat64 `json:"spent_quantity"`
	SpentTotal    float64         `json:"spent_total"`

	RemainingQuantity sql.NullFloat64 `json:"remaining_quantity"`
	RemainingCost     sql.NullFloat64 `json:"remaining_cost"`
	RemainingTotal    float64         `json:"remaining_total"`

	UpdatedBudget float64 `json:"updated_budget"`

	CompanyId uuid.UUID `json:"company_id"`
}

type ProjectData struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	NetArea   float64   `json:"net_area"`
	GrossArea float64   `json:"gross_area"`
}

type BudgetItemData struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Level      uint8     `json:"level"`
	Accumulate bool      `json:"accumulate"`
}
