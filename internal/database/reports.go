package database

import (
	"bca-go-final/internal/types"
	"log"
	"time"

	"github.com/google/uuid"
)

func (s *service) GetBalance(companyId, projectId uuid.UUID, date time.Time) types.BalanceResponse {
	invoices := []types.InvoiceResponse{}
	query := `
	select id, supplier_id, supplier_number, supplier_name, supplier_contact_name, supplier_contact_email, 
	supplier_contact_phone, project_id, project_name, project_is_active, invoice_number, invoice_date, invoice_total,
	company_id
	from vw_invoice
	where extract(year from invoice_date) = $1 and extract(month from invoice_date) = $2 and company_id = $3
	and project_id = $4
	order by invoice_date desc, supplier_name, invoice_number
	`
	rows, err := s.db.Query(query, date.Year(), date.Month(), companyId, projectId)
	if err != nil {
		log.Fatal("Error in select", err)
		return types.BalanceResponse{}
	}
	defer rows.Close()
	var total float64 = 0

	for rows.Next() {
		i := types.InvoiceResponse{}

		if err := rows.Scan(
			&i.Id,
			&i.Supplier.ID,
			&i.Supplier.SupplierId,
			&i.Supplier.Name,
			&i.Supplier.ContactName,
			&i.Supplier.ContactEmail,
			&i.Supplier.ContactPhone,
			&i.Project.ID,
			&i.Project.Name,
			&i.Project.IsActive,
			&i.InvoiceNumber,
			&i.InvoiceDate,
			&i.InvoiceTotal,
			&i.CompanyId,
		); err != nil {
			log.Println("Error in scan", err)
			return types.BalanceResponse{}
		}
		invoices = append(invoices, i)
		total += i.InvoiceTotal
	}

	return types.BalanceResponse{Invoices: invoices, Total: total}
}
