package types

import "github.com/google/uuid"

type User struct {
	Id        uuid.UUID
	Email     string
	CompanyId uuid.UUID
	Name      string
	RoleId    string
}

type UserCreate struct {
	Id        uuid.UUID
	Email     string
	Password  string
	CompanyId uuid.UUID
	Name      string
	RoleId    string
}
