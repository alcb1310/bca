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
	ProjectId    uuid.UUID
	BudgetItemId uuid.UUID

	Quantity float64
	Cost     float64

	CompanyId uuid.UUID
}

// Structure for reading budgets
type GetBudget struct {
	Project    project
	BudgetItem budgetItem

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

type project struct {
	ID        uuid.UUID
	Name      string
	NetArea   float64
	GrossArea float64
}

type budgetItem struct {
	ID         uuid.UUID
	Code       string
	Name       string
	Level      uint8
	Accumulate bool
}
