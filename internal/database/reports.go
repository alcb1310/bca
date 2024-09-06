package database

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/alcb1310/bca/internal/types"
)

func (s *service) GetBalance(companyId, projectId uuid.UUID, date time.Time) types.BalanceResponse {
	invoices := []types.InvoiceResponse{}
	query := `
	select id, supplier_id, supplier_number, supplier_name, supplier_contact_name, supplier_contact_email, 
	supplier_contact_phone, project_id, project_name, project_is_active, invoice_number, invoice_date, invoice_total,
	company_id, is_balanced
	from vw_invoice
	where extract(year from invoice_date) = $1 and extract(month from invoice_date) = $2 and company_id = $3
	and project_id = $4
	order by invoice_date desc, supplier_name, invoice_number
	`
	rows, err := s.db.Query(query, date.Year(), date.Month(), companyId, projectId)
	if err != nil {
		slog.Error("Error in select", "err", err)
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
			&i.IsBalanced,
		); err != nil {
			slog.Error("Error in select", "err", err)
			return types.BalanceResponse{}
		}
		invoices = append(invoices, i)
		total += i.InvoiceTotal
	}

	return types.BalanceResponse{Invoices: invoices, Total: total}
}

func (s *service) GetHistoricByProject(companyId, projectId uuid.UUID, date time.Time, level uint8) []types.GetBudget {
	var rows *sql.Rows
	var err error
	query := `
        SELECT
            project_id, project_name, project_gross_area, project_net_area,
            budget_item_id, budget_item_code, budget_item_name, budget_item_level, budget_item_accumulate,
            initial_quantity, initial_cost, initial_total,
            spent_quantity, spent_total,
            remaining_quantity, remaining_cost, remaining_total,
            updated_budget, company_id
        FROM vw_historic
        WHERE company_id = $1 and project_id = $2 AND budget_item_level <= $3 and
	        extract(year from date) = $4 and extract(month from date) = $5
		ORDER BY budget_item_code
		`
	rows, err = s.db.Query(query, companyId, projectId, level, date.Year(), date.Month())
	if err != nil {
		slog.Error("Error in query", "err", err)
		return nil
	}
	defer rows.Close()
	budgets := []types.GetBudget{}

	for rows.Next() {
		b := types.GetBudget{}
		if err := rows.Scan(
			&b.Project.ID, &b.Project.Name, &b.Project.GrossArea, &b.Project.NetArea,
			&b.BudgetItem.ID, &b.BudgetItem.Code, &b.BudgetItem.Name, &b.BudgetItem.Level, &b.BudgetItem.Accumulate,
			&b.InitialQuantity, &b.InitialCost, &b.InitialTotal,
			&b.SpentQuantity, &b.SpentTotal,
			&b.RemainingQuantity, &b.RemainingCost, &b.RemainingTotal,
			&b.UpdatedBudget, &b.CompanyId,
		); err != nil {
			slog.Error("Error in scan", "err", err)
			return nil
		}
		budgets = append(budgets, b)
	}

	return budgets
}

func (s *service) GetSpentByBudgetItem(companyId, projectId, budgetItemId uuid.UUID, date time.Time, ids []uuid.UUID) float64 {
	query := `
	    select sum(total)
		from vw_invoice_details where company_id=$1 and extract(year from invoice_date)=$2 and
		extract(month from invoice_date)=$3 and project_id=$4 and budget_item_id=any($5)
	`
	var total *float64
	s.db.QueryRow(query, companyId, date.Year(), date.Month(), projectId, pq.Array(ids)).Scan(&total)
	if total == nil {
		return 0
	}

	return *total
}

func (s *service) GetDetailsByBudgetItem(companyId, projectId, budgetItemId uuid.UUID, date time.Time, ids []uuid.UUID) []types.InvoiceDetails {
	query := `
	    select invoice_id, invoice_number, invoice_total, invoice_date, project_id, project_name, supplier_id, supplier_number,
		supplier_name, budget_item_id, budget_item_name, budget_item_code, budget_item_level, quantity, cost, total, company_id
		from vw_invoice_details where company_id=$1 and extract(year from invoice_date)=$2 and
		extract(month from invoice_date)=$3 and project_id=$4 and budget_item_id=any($5)
	`
	row, err := s.db.Query(query, companyId, date.Year(), date.Month(), projectId, pq.Array(ids))
	if err != nil {
		slog.Error("Error in query", "err", err)
		return []types.InvoiceDetails{}
	}
	defer row.Close()

	returnInvoiceDetails := []types.InvoiceDetails{}

	for row.Next() {
		i := types.InvoiceDetails{}

		if err := row.Scan(
			&i.InvoiceId, &i.InvoiceNumber, &i.InvoiceTotal, &i.InvoiceDate, &i.ProjectId, &i.ProjectName, &i.SupplierId, &i.SupplierNumber,
			&i.SupplierName, &i.BudgetItemId, &i.BudgetItemName, &i.BudgetItemCode, &i.BudgetItemLevel, &i.Quantity, &i.Cost, &i.Total, &i.CompanyId,
		); err != nil {
			slog.Error("Error in scan", "err", err)
			return []types.InvoiceDetails{}
		}
		returnInvoiceDetails = append(returnInvoiceDetails, i)

	}

	return returnInvoiceDetails
}
