package types

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceDetailsResponse struct {
	Id             uuid.UUID `json:"id"`
	BudgetItemId   uuid.UUID `json:"budget_item_id"`
	BudgetItemCode string    `json:"budget_item_code"`
	BudgetItemName string    `json:"budget_item_name"`
	Quantity       float64   `json:"quantity"`
	Cost           float64   `json:"cost"`
	Total          float64   `json:"total"`
	InvoiceTotal   float64   `json:"invoice_total"`
	CompanyId      uuid.UUID `json:"company_id"`
}

type InvoiceDetailCreate struct {
	InvoiceId             uuid.UUID
	BudgetItemId          uuid.UUID
	Quantity, Cost, Total float64
	CompanyId             uuid.UUID
}

type InvoiceDetails struct {
	InvoiceId       uuid.UUID
	InvoiceDate     time.Time
	InvoiceNumber   string
	InvoiceTotal    float64
	ProjectId       uuid.UUID
	ProjectName     string
	SupplierId      uuid.UUID
	SupplierName    string
	SupplierNumber  string
	BudgetItemId    uuid.UUID
	BudgetItemCode  string
	BudgetItemName  string
	BudgetItemLevel uint8
	Quantity        float64
	Cost            float64
	Total           float64
	CompanyId       uuid.UUID
}
