package database

import "bca-go-final/internal/types"

func (s *service) Register(company *types.Company) error {
	sql := "insert into company (ruc, name, employees, is_active) values ($1, $2, $3, $4)"
	_, err := s.db.Exec(sql, company.ID, company.Ruc, company.Name, company.Employees, company.IsActive)
	if err != nil {
		return err
	}

	return nil
}
