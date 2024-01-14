package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetAllDetails(invoiceId, companyId uuid.UUID) ([]types.InvoiceDetailsResponse, error) {
	details := []types.InvoiceDetailsResponse{}
	query := "select invoice_id, budget_item_id, budget_item_code, budget_item_name, quantity, cost, total, invoice_total  from vw_invoice_details where invoice_id = $1 and company_id = $2"

	rows, err := s.db.Query(query, invoiceId, companyId)
	if err != nil {
		return []types.InvoiceDetailsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var detail types.InvoiceDetailsResponse
		if err := rows.Scan(&detail.Id, &detail.BudgetItemId, &detail.BudgetItemCode, &detail.BudgetItemName, &detail.Quantity, &detail.Cost, &detail.Total, &detail.InvoiceTotal); err != nil {
			return []types.InvoiceDetailsResponse{}, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (s *service) AddDetail(detail types.InvoiceDetailCreate) error {
	// TODO: Start a transaction
	// TODO: On Error rollback changes

	// TODO: Add detail to the details table
	// TODO: In the invoice table update the total
	// TODO: Update the budget for the saved budget budget item
	// TODO: Update the budget for the parent budget item

	// TODO: On Success commit changes
	return nil
}
