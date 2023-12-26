package types

import "github.com/google/uuid"

type Supplier struct {
	ID uuid.UUID `json:"id"`

	SupplierId string `json:"supplier_id"`
	Name       string `json:"name"`

	ContactName  *string `json:"contact_name"`
	ContactEmail *string `json:"contact_email"`
	ContactPhone *string `json:"contact_phone"`

	CompanyId uuid.UUID `json:"company_id"`
}
