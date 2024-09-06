package types

import (
	"time"

	"github.com/google/uuid"
)

type ContextPayload struct {
	Id         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	CompanyId  uuid.UUID `json:"company_id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	IsLoggedIn bool      `json:"is_logged_in"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}
