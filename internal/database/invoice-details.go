package database

import (
	"errors"
	"log"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
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
	query := "select is_balanced from invoice where id = $1 and company_id = $2"
	var isBalanced bool
	if err := s.db.QueryRow(query, detail.InvoiceId, detail.CompanyId).Scan(&isBalanced); err != nil {
		return err
	}
	if isBalanced {
		return errors.New("La factura ya se encuentra balanceada")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query = "insert into invoice_details (invoice_id, budget_item_id, quantity, cost, total, company_id) values ($1, $2, $3, $4, $5, $6)"
	if _, err := tx.Exec(query, detail.InvoiceId, detail.BudgetItemId, detail.Quantity, detail.Cost, detail.Total, detail.CompanyId); err != nil {
		return err
	}
	query = "update invoice set invoice_total = invoice_total + $1 where id = $2 and company_id = $3"
	if _, err := tx.Exec(query, detail.Total, detail.InvoiceId, detail.CompanyId); err != nil {
		return err
	}

	var projectId uuid.UUID
	query = "select project_id from invoice where id = $1 and company_id = $2"
	if err := tx.QueryRow(query, detail.InvoiceId, detail.CompanyId).Scan(&projectId); err != nil {
		return err
	}

	query = "select spent_quantity, spent_total, remaining_quantity, remaining_cost, remaining_total, updated_budget from budget where project_id = $1 and budget_item_id = $2 and company_id = $3"
	var spentQuantity, spentTotal, remainingQuantity, remainingCost, remainingTotal, updatedBudget float64
	if err := tx.QueryRow(query, projectId, detail.BudgetItemId, detail.CompanyId).Scan(&spentQuantity, &spentTotal, &remainingQuantity, &remainingCost, &remainingTotal, &updatedBudget); err != nil {
		return err
	}

	newToSpendTotal := (remainingQuantity - detail.Quantity) * detail.Cost
	newSpentTotal := spentTotal + detail.Total
	newUpdatedBudget := newSpentTotal + newToSpendTotal

	query = `
		update budget set spent_quantity = spent_quantity + $1, spent_total = spent_total + $3,
		remaining_quantity = remaining_quantity - $1, remaining_cost = $2, remaining_total = $4,
		updated_budget = $5
		where project_id = $6 and budget_item_id = $7 and company_id = $8
	`
	if _, err := tx.Exec(query, detail.Quantity, detail.Cost, detail.Total, newToSpendTotal, newUpdatedBudget, projectId, detail.BudgetItemId, detail.CompanyId); err != nil {
		return err
	}

	updatedDiff := newUpdatedBudget - updatedBudget
	remainingDiff := newToSpendTotal - remainingTotal

	var parentId *uuid.UUID
	parentId = &detail.BudgetItemId
	for true {
		query = "select parent_id from budget_item where id = $1 and company_id = $2"
		if err := tx.QueryRow(query, parentId, detail.CompanyId).Scan(&parentId); err != nil {
			return err
		}
		if parentId == nil || parentId == &uuid.Nil {
			break
		}

		query = `
		update budget set spent_total = spent_total + $1, remaining_total = remaining_total + $2,
		updated_budget = updated_budget + $3
		where project_id = $4 and budget_item_id = $5 and company_id = $6
		`
		if _, err := tx.Exec(query, detail.Total, remainingDiff, updatedDiff, projectId, parentId, detail.CompanyId); err != nil {
			return err
		}

	}

	tx.Commit()
	return nil
}

func (s *service) DeleteDetail(invoiceId, budgetItemId, companyId uuid.UUID) error {
	var quantity, cost, total float64
	query := "select quantity, cost, total from invoice_details where invoice_id = $1 and budget_item_id = $2 and company_id = $3"
	if err := s.db.QueryRow(query, invoiceId, budgetItemId, companyId).Scan(&quantity, &cost, &total); err != nil {
		log.Println("Error en el select: ", err)
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query = "delete from invoice_details where invoice_id = $1 and budget_item_id = $2 and company_id = $3"
	if _, err := tx.Exec(query, invoiceId, budgetItemId, companyId); err != nil {
		log.Println("Error en el delete: ", err)
		return err
	}
	query = "update invoice set invoice_total = invoice_total - $1 where id = $2 and company_id = $3"
	if _, err := tx.Exec(query, total, invoiceId, companyId); err != nil {
		log.Println("Error en el update: ", err)
		return err
	}
	var projectId uuid.UUID
	query = "select project_id from invoice where id = $1 and company_id = $2"
	if err := tx.QueryRow(query, invoiceId, companyId).Scan(&projectId); err != nil {
		return err
	}

	query = `
	update budget set spent_quantity = spent_quantity - $1, spent_total = spent_total - $2,
		remaining_quantity = remaining_quantity + $1, remaining_cost = $3, remaining_total = remaining_total + $2
	where project_id = $4 and budget_item_id = $5 and company_id = $6
	`
	if _, err := tx.Exec(query, quantity, total, cost, projectId, budgetItemId, companyId); err != nil {
		log.Println("Error en el update budget: ", err)
		return err
	}

	var parentId *uuid.UUID
	parentId = &budgetItemId
	for true {
		query = "select parent_id from budget_item where id = $1 and company_id = $2"
		if err := tx.QueryRow(query, parentId, companyId).Scan(&parentId); err != nil {
			return err
		}
		if parentId == nil || parentId == &uuid.Nil {
			break
		}
		query = `
		update budget set spent_total = spent_total - $1, remaining_total = remaining_total + $1
		where project_id = $2 and budget_item_id = $3 and company_id = $4
		`
		if _, err := tx.Exec(query, total, projectId, parentId, companyId); err != nil {
			log.Println("Error en el update budget: ", err)
			return err
		}

	}

	tx.Commit()
	return nil
}
