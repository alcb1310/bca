package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetBudgetItems(companyId uuid.UUID) ([]types.BudgetItem, error) {
	sql := "select id, code, name, level, accumulate, parent_id, company_id from vw_budget_item where company_id = $1"

	rows, err := s.db.Query(sql, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bis := []types.BudgetItem{}

	for rows.Next() {
		bi := types.BudgetItem{}
		if err := rows.Scan(&bi.ID, &bi.Code, &bi.Name, &bi.Level, &bi.Accumulate, &bi.ParentId, &bi.CompanyId); err != nil {
			return nil, err
		}
		bis = append(bis, bi)
	}

	return bis, nil
}

func (s *service) CreateBudgetItem(bi *types.BudgetItem) error {
	var level uint8 = 1
	if bi.ParentId != nil {
		sql := "select level from budget_item where id = $1 and company_id = $2"
		err := s.db.QueryRow(sql, bi.ParentId, bi.CompanyId).Scan(&level)
		if err != nil {
			return err
		}
		level++
	}
	bi.Level = level

	sql := "insert into budget_item (code, name, level, accumulate, parent_id, company_id) values ($1, $2, $3, $4, $5, $6) returning id"
	err := s.db.QueryRow(sql, bi.Code, bi.Name, bi.Level, bi.Accumulate, bi.ParentId, bi.CompanyId).Scan(&bi.ID)
	return err
}
