package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
)

func (s *service) CreateClosure(companyId, projectId uuid.UUID, date time.Time) error {
	log.Println(fmt.Sprintf("Procesando el cierre del proyecto: %s para la fecha: %s", projectId, utils.ConvertDate(date)))
	var existingDate time.Time

	query := "select date from historic where company_id = $1 and project_id = $2 and extract(year from date) = $3 and extract(month from date) = $4"
	row := s.db.QueryRow(query, companyId, projectId, date.Year(), date.Month())
	row.Scan(&existingDate)
	if existingDate != (time.Time{}) {
		return errors.New(fmt.Sprintf("Ya existe un cierre para el proyecto: %s para la fecha: %s", projectId, utils.ConvertDate(date)))
	}

	query = `
	    select project_id, budget_item_id, initial_quantity, initial_cost, initial_total, spent_quantity, spent_total,
		remaining_quantity, remaining_cost, remaining_total, updated_budget, company_id
		from budget 
		where project_id = $1 and company_id = $2
	`
	rows, err := s.db.Query(query, projectId, companyId)
	if err != nil {
		log.Println("Error en el Query")
		return err
	}
	defer rows.Close()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for rows.Next() {
		b := types.Budget{}
		if err := rows.Scan(&b.ProjectId, &b.BudgetItemId, &b.InitialQuantity, &b.InitialCost, &b.InitialTotal, &b.SpentQuantity, &b.SpentTotal,
			&b.RemainingQuantity, &b.RemainingCost, &b.RemainingTotal, &b.UpdatedBudget, &b.CompanyId); err != nil {
			log.Println("Error en el Scan")
			return err
		}

		query = `
			insert into historic
			(project_id, budget_item_id, initial_quantity, initial_cost, initial_total, spent_quantity, spent_total,
			remaining_quantity, remaining_cost, remaining_total, updated_budget, company_id, date)
			values
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`
		if _, err := tx.Exec(query, b.ProjectId, b.BudgetItemId, b.InitialQuantity, b.InitialCost, b.InitialTotal, b.SpentQuantity, b.SpentTotal,
			b.RemainingQuantity, b.RemainingCost, b.RemainingTotal, b.UpdatedBudget, b.CompanyId, date); err != nil {
			log.Println("Error en el insert")
			return err
		}
	}

	query = "update invoice set is_balanced = true where company_id = $1 and project_id = $2 and extract(year from invoice_date) = $3 and extract(month from invoice_date) = $4"
	if _, err := tx.Exec(query, companyId, projectId, date.Year(), date.Month()); err != nil {
		log.Println("Error en el update")
		return err
	}

	query = "update project set last_closure = $1 where id = $2 and company_id = $3"
	if _, err := tx.Exec(query, date, projectId, companyId); err != nil {
		log.Println("Error en el update project")
		return err
	}

	tx.Commit()
	log.Println(fmt.Sprintf("Terminado el cierre del proyecto: %s para la fecha: %s", projectId, utils.ConvertDate(date)))
	return nil
}
