package types

import (
	"github.com/google/uuid"
	"time"
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
	InvoiceId    uuid.UUID `json:"invoiceId"`
	BudgetItemId uuid.UUID `json:"budget_item_id"`
	Quantity     float64   `json:"quantity"`
	Cost         float64   `json:"cost"`
	Total        float64   `json:"total"`
	CompanyId    uuid.UUID `json:"companyId"`
}

type InvoiceDetails struct {
	InvoiceId       uuid.UUID `json:"invoice_id"`
	InvoiceDate     time.Time `json:"invoice_date"`
	InvoiceNumber   string    `json:"invoice_number"`
	InvoiceTotal    float64   `json:"invoice_total"`
	ProjectId       uuid.UUID `json:"project_id"`
	ProjectName     string    `json:"project_name"`
	SupplierId      uuid.UUID `json:"supplier_id"`
	SupplierName    string    `json:"supplier_name"`
	SupplierNumber  string    `json:"supplier_number"`
	BudgetItemId    uuid.UUID `json:"budget_item_id"`
	BudgetItemCode  string    `json:"budget_item_code"`
	BudgetItemName  string    `json:"budget_item_name"`
	BudgetItemLevel uint8     `json:"budgget_item_level"`
	Quantity        float64   `json:"quantity"`
	Cost            float64   `json:"cost"`
	Total           float64   `json:"total"`
	CompanyId       uuid.UUID `json:"company_id"`
}
