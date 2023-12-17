package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetAllSuppliers(companyId uuid.UUID) ([]types.Supplier, error) {
	suppliers := []types.Supplier{}

	sql := "SELECT id, supplier_id, name, contact_name, contact_email, contact_phone, company_id FROM supplier WHERE company_id = $1"
	rows, err := s.db.Query(sql, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := types.Supplier{}
		if err := rows.Scan(&s.ID, &s.SupplierId, &s.Name, &s.ContactName, &s.ContactEmail, &s.ContactPhone, &s.CompanyId); err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}

	return suppliers, nil
}
