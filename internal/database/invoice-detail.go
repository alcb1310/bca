package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetInvoiceDetails(invoiceId, companyId uuid.UUID) ([]types.InvoiceDetailResponse, error) {
	id := []types.InvoiceDetailResponse{}
	query := `
	select invoice_id, supplier_id, supplier_number, supplier_name, supplier_contact_name, supplier_contact_email, supplier_contact_phone,
		project_id, project_name, project_is_active, invoice_number, invoice_date, invoice_total,
		budget_item_id, budget_item_code, budget_item_name, budget_item_level, budget_item_accumulate, budget_item_parent_id,
		quantity, cost, total, company_id
	from vw_invoice_details
	where company_id = $1 and invoice_id = $2
	`

	rows, err := s.db.Query(query, companyId, invoiceId)
	if err != nil {
		return id, err
	}
	defer rows.Close()

	for rows.Next() {
		i := types.InvoiceDetailResponse{}
		if err := rows.Scan(
			&i.Invoice.Id, &i.Invoice.Supplier.ID, &i.Invoice.Supplier.SupplierId, &i.Invoice.Supplier.Name, &i.Invoice.Supplier.ContactName, &i.Invoice.Supplier.ContactEmail, &i.Invoice.Supplier.ContactPhone,
			&i.Invoice.Project.ID, &i.Invoice.Project.Name, &i.Invoice.Project.IsActive, &i.Invoice.InvoiceNumber, &i.Invoice.InvoiceDate, &i.Invoice.InvoiceTotal,
			&i.BudgetItem.ID, &i.BudgetItem.Code, &i.BudgetItem.Name, &i.BudgetItem.Level, &i.BudgetItem.Accumulate, &i.BudgetItem.ParentId,
			&i.Quantity, &i.Cost, &i.Total, &i.CompanyId); err != nil {
			return id, err
		}

		i.Invoice.CompanyId = companyId
		i.Invoice.Supplier.CompanyId = companyId
		i.Invoice.Project.CompanyId = companyId
		i.BudgetItem.CompanyId = companyId

		id = append(id, i)
	}

	return id, nil
}

func (s *service) CreateInvoiceDetail(detail *types.InvoiceDetail) (types.InvoiceDetail, error) {
	// TODO: implement create invoice details
	return types.InvoiceDetail{}, nil
}

func (s *service) UpdateInvoiceDetail(detail *types.InvoiceDetail) error {
	// TODO: implement update invoice details
	return nil
}

func (s *service) DeleteInvoiceDetail(detailId, companyId uuid.UUID) error {
	// TODO: implement delete invoice details
	return nil
}

func (s *service) GetOneInvoiceDetail(invoiceId, detailId, companyId uuid.UUID) (types.InvoiceDetailResponse, error) {
	// TODO: implement get invoice details
	return types.InvoiceDetailResponse{}, nil
}
