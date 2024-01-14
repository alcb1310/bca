package types

import "github.com/google/uuid"

type InvoiceDetailsResponse struct {
	Id             uuid.UUID
	BudgetItemId   uuid.UUID
	BudgetItemCode string
	BudgetItemName string
	Quantity       float64
	Cost           float64
	Total          float64
	InvoiceTotal   float64
	CompanyId      uuid.UUID
}

type InvoiceDetailCreate struct {
	InvoiceId             uuid.UUID
	BudgetItemId          uuid.UUID
	Quantity, Cost, Total float64
	CompanyId             uuid.UUID
}
