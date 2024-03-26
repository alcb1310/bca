package types

import "github.com/google/uuid"

type Category struct {
	Id        uuid.UUID
	Name      string
	CompanyId uuid.UUID
}
