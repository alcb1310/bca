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
	CompanyId     uuid.UUID
}

type InvoiceResponse struct {
	Id            uuid.UUID
	Supplier      Supplier
	Project       Project
	InvoiceNumber string
	InvoiceDate   time.Time
	InvoiceTotal  float64
	IsBalanced    bool
	CompanyId     uuid.UUID
}
