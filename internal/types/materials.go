package types

import "github.com/google/uuid"

type Material struct {
	Id        uuid.UUID
	Code      string
	Name      string
	Unit      string
	Category  Category
	CompanyId uuid.UUID
}
