package types

import "github.com/google/uuid"

type BalanceResponse struct {
	Invoices []InvoiceResponse
	Total    float64
}

type Spent struct {
	Spent      float64
	BudgetItem BudgetItem
}

type SpentResponse struct {
	Spent   []Spent
	Total   float64
	Project uuid.UUID
}
