package types

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceCreate struct {
	Id            *uuid.UUID
	SupplierId    *uuid.UUID
	ProjectId     *uuid.UUID
	InvoiceNumber *string
	InvoiceDate   *time.Time
	IsBalanced    bool
	CompanyId     uuid.UUID
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
