package types

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceCreate struct {
	Id            *uuid.UUID `json:"id"`
	SupplierId    *uuid.UUID `json:"supplier_id"`
	ProjectId     *uuid.UUID `json:"project_id"`
	InvoiceNumber *string    `json:"invoice_number"`
	InvoiceDate   *time.Time `json:"invoice_date"`
	InvoiceTotal  float64    `json:"invoice_total"`
	IsBalanced    bool       `json:"is_balanced"`
	CompanyId     uuid.UUID  `json:"company_id"`
}

type InvoiceResponse struct {
	Id            uuid.UUID `json:"id"`
	Supplier      Supplier  `json:"supplier"`
	Project       Project   `json:"project"`
	InvoiceNumber string    `json:"invoice_number"`
	InvoiceDate   time.Time `json:"invoice_date"`
	InvoiceTotal  float64   `json:"invoice_total"`
	IsBalanced    bool      `json:"is_balanced"`
	CompanyId     uuid.UUID `json:"company_id"`
}
