package types

import "github.com/google/uuid"

type Supplier struct {
	ID uuid.UUID

	SupplierId string
	Name       string

	ContactName  *string
	ContactEmail *string
	ContactPhone *string

	CompanyId uuid.UUID
}
