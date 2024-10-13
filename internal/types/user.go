package types

import "github.com/google/uuid"

type User struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CompanyId uuid.UUID `json:"company_id"`
	Name      string    `json:"name"`
	RoleId    string    `json:"role_id"`
}

type UserCreate struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CompanyId uuid.UUID `json:"company_id"`
	Name      string    `json:"name"`
	RoleId    string    `json:"role_id"`
}
