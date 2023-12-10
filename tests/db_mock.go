package tests

import "bca-go-final/internal/types"

type DBMock struct{}

func (s *DBMock) Login(l *types.Login) (string, error) {
	return "", nil
}

func (s *DBMock) CreateCompany(company *types.CompanyCreate) error {
	return nil
}

func (s *DBMock) Health() map[string]string {
	health := make(map[string]string)
	health["message"] = "It's healthy"
	return health
}
