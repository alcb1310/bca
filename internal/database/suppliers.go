package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetAllSuppliers(companyId uuid.UUID, search string) ([]types.Supplier, error) {
	suppliers := []types.Supplier{}

	se := "%" + search + "%"
	sql := "SELECT id, supplier_id, name, contact_name, contact_email, contact_phone, company_id FROM supplier WHERE company_id = $1 and (supplier_id like $2 or name like $2) ORDER BY name"
	rows, err := s.db.Query(sql, companyId, se)
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

func (s *service) CreateSupplier(supplier *types.Supplier) error {
	sql := "INSERT INTO supplier (supplier_id, name, contact_name, contact_email, contact_phone, company_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err := s.db.QueryRow(sql, supplier.SupplierId, supplier.Name, supplier.ContactName, supplier.ContactEmail, supplier.ContactPhone, supplier.CompanyId).Scan(&supplier.ID)
	return err
}

func (s *service) GetOneSupplier(id, companyId uuid.UUID) (types.Supplier, error) {
	supplier := types.Supplier{}
	sql := "SELECT id, supplier_id, name, contact_name, contact_email, contact_phone, company_id FROM supplier WHERE id = $1 AND company_id = $2"
	err := s.db.QueryRow(sql, id, companyId).Scan(&supplier.ID, &supplier.SupplierId, &supplier.Name, &supplier.ContactName, &supplier.ContactEmail, &supplier.ContactPhone, &supplier.CompanyId)
	return supplier, err
}

func (s *service) UpdateSupplier(supplier *types.Supplier) error {
	sql := "UPDATE supplier SET supplier_id = $1, name = $2, contact_name = $3, contact_email = $4, contact_phone = $5 WHERE id = $6 AND company_id = $7"
	_, err := s.db.Exec(sql, supplier.SupplierId, supplier.Name, supplier.ContactName, supplier.ContactEmail, supplier.ContactPhone, supplier.ID, supplier.CompanyId)
	return err
}
