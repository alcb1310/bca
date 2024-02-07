package types

import (
	"database/sql"

	"github.com/google/uuid"
)

type Supplier struct {
	ID uuid.UUID

	SupplierId string
	Name       string

	ContactName  sql.NullString
	ContactEmail sql.NullString
	ContactPhone sql.NullString

	CompanyId uuid.UUID
}
