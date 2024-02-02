package types

import "github.com/google/uuid"

type Budget struct {
	ProjectId    uuid.UUID
	BudgetItemId uuid.UUID

	InitialQuantity *float64
	InitialCost     *float64
	InitialTotal    float64

	SpentQuantity *float64
	SpentTotal    float64

	RemainingQuantity *float64
	RemainingCost     *float64
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

	InitialQuantity *float64
	InitialCost     *float64
	InitialTotal    float64

	SpentQuantity *float64
	SpentTotal    float64

	RemainingQuantity *float64
	RemainingCost     *float64
	RemainingTotal    float64

	UpdatedBudget float64

	CompanyId uuid.UUID
}

type project struct {
	ID   uuid.UUID
	Name string
}

type budgetItem struct {
	ID         uuid.UUID
	Code       string
	Name       string
	Level      uint8
	Accumulate bool
}
