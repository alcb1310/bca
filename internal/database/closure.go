package database

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *service) CreateClosure(companyId, projectId uuid.UUID, date time.Time) error {
	slog.Info("Procesando el cierre del proyecto", "project_id", projectId, "date", utils.ConvertDate(date))
	var existingDate time.Time

	query := "select date from historic where company_id = $1 and project_id = $2 and extract(year from date) = $3 and extract(month from date) = $4"
	row := s.db.QueryRow(query, companyId, projectId, date.Year(), date.Month())
	row.Scan(&existingDate)
	if existingDate != (time.Time{}) {
		query := "select name from project where id = $1 and company_id = $2"
		var projectName string
		s.db.QueryRow(query, projectId, companyId).Scan(&projectName)
		return errors.New(fmt.Sprintf("Ya existe un cierre para el proyecto: %s para la fecha: %s", projectName, utils.ConvertDate(date)))
	}

	query = `
	    select project_id, budget_item_id, initial_quantity, initial_cost, initial_total, spent_quantity, spent_total,
		remaining_quantity, remaining_cost, remaining_total, updated_budget, company_id
		from budget 
		where project_id = $1 and company_id = $2
	`
	rows, err := s.db.Query(query, projectId, companyId)
	if err != nil {
		slog.Error("Error en el Query", "query", query, "err", err)
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
			slog.Error("Error en el Scan", "err", err)
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
			slog.Error("Error en el insert", "err", err)
			return err
		}
	}

	query = "update invoice set is_balanced = true where company_id = $1 and project_id = $2 and extract(year from invoice_date) = $3 and extract(month from invoice_date) = $4"
	if _, err := tx.Exec(query, companyId, projectId, date.Year(), date.Month()); err != nil {
		slog.Error("Error en el update", "err", err)
		return err
	}

	query = "update project set last_closure = $1 where id = $2 and company_id = $3"
	if _, err := tx.Exec(query, date, projectId, companyId); err != nil {
		slog.Error("Error en el update project", "err", err)
		return err
	}

	tx.Commit()
	slog.Info("Terminado el cierre del proyecto", "project_id", projectId, "date", utils.ConvertDate(date))
	return nil
}
