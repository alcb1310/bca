package database

import "bca-go-final/internal/types"

func (s *service) CreateCompany(company *types.Company) error {
	sql := "insert into company (ruc, name, employees) values ($1, $2, $3)"
	_, err := s.db.Exec(sql, company.Ruc, company.Name, company.Employees)
	if err != nil {
		return err
	}

	return nil
}
