package types

import "github.com/google/uuid"

type Budget struct {
	ProjectId    uuid.UUID `json:"project_id"`
	BudgetItemId uuid.UUID `json:"budget_item_id"`

	InitialQuantity *float64 `json:"initial_quantity"`
	InitialCost     *float64 `json:"initial_cost"`
	InitialTotal    float64  `json:"initial_total"`

	SpentQuantity *float64 `json:"spent_quantity"`
	SpentTotal    float64  `json:"spent_total"`

	RemainingQuantity *float64 `json:"remaining_quantity"`
	RemainingCost     *float64 `json:"remaining_cost"`
	RemainingTotal    float64  `json:"remaining_total"`

	UpdatedBudget float64 `json:"updated_budget"`

	CompanyId uuid.UUID `json:"company_id"`
}

// Structure for creating budget
type CreateBudget struct {
	ProjectId    uuid.UUID `json:"project_id"`
	BudgetItemId uuid.UUID `json:"budget_item_id"`

	Quantity float64 `json:"quantity"`
	Cost     float64 `json:"cost"`

	CompanyId uuid.UUID `json:"company_id"`
}

// Structure for reading budgets
type GetBudget struct {
	Project    project    `json:"project"`
	BudgetItem budgetItem `json:"budget_item"`

	InitialQuantity *float64 `json:"initial_quantity"`
	InitialCost     *float64 `json:"initial_cost"`
	InitialTotal    float64  `json:"initial_total"`

	SpentQuantity *float64 `json:"spent_quantity"`
	SpentTotal    float64  `json:"spent_total"`

	RemainingQuantity *float64 `json:"remaining_quantity"`
	RemainingCost     *float64 `json:"remaining_cost"`
	RemainingTotal    float64  `json:"remaining_total"`

	UpdatedBudget float64 `json:"updated_budget"`

	CompanyId uuid.UUID `json:"company_id"`
}

type project struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type budgetItem struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Level      uint8     `json:"level"`
	Accumulate bool      `json:"accumulate"`
}
