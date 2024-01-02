package types

import "github.com/google/uuid"

type InvoiceDetail struct {
	// TODO: implement invoice details structure
}

type InvoiceDetailResponse struct {
	// TODO: implement invoice details response structure
	Invoice    InvoiceResponse `json:"invoice"`
	BudgetItem BudgetItem      `json:"budget_item"`
	Quantity   float64         `json:"quantity"`
	Cost       float64         `json:"cost"`
	Total      float64         `json:"total"`
	CompanyId  uuid.UUID       `json:"company_id"`
}
