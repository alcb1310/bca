package types

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CompanyID uuid.UUID
}

type CreateUser struct {
	User
	Password string
}
