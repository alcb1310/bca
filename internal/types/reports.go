package types

import "github.com/google/uuid"

type BalanceResponse struct {
	Invoices []InvoiceResponse `json:"invoices"`
	Total    float64           `json:"total"`
}

type Spent struct {
	Spent      float64    `json:"spent"`
	BudgetItem BudgetItem `json:"budget_item"`
}

type SpentResponse struct {
	Spent   []Spent   `json:"spent"`
	Total   float64   `json:"total"`
	Project uuid.UUID `json:"project"`
}
