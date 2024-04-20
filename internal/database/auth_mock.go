package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) CreateCompany(company *types.CompanyCreate) error {
	return nil
}

func (s ServiceMock) Login(l *types.Login) (string, error) {
	return "", nil
}

func (s ServiceMock) RegenerateToken(token string, user uuid.UUID) error {
	return nil
}

func (s ServiceMock) IsLoggedIn(token string, user uuid.UUID) bool {
	return true
}
