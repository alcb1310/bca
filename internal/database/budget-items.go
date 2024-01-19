package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetBudgetItems(companyId uuid.UUID, search string) ([]types.BudgetItemResponse, error) {
	sql := "select id, code, name, level, accumulate, parent_id, parent_code, parent_name, company_id from vw_budget_item where company_id = $1 and (name like $2 or code like $2) order by code"

	searchTerm := "%" + search + "%"
	rows, err := s.db.Query(sql, companyId, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bis := []types.BudgetItemResponse{}

	for rows.Next() {
		bi := types.BudgetItemResponse{}
		if err := rows.Scan(&bi.ID, &bi.Code, &bi.Name, &bi.Level, &bi.Accumulate, &bi.ParentId, &bi.ParentCode, &bi.ParentName, &bi.CompanyId); err != nil {
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

func (s *service) GetOneBudgetItem(id uuid.UUID, companyId uuid.UUID) (*types.BudgetItem, error) {
	bi := &types.BudgetItem{}
	sql := "select id, code, name, level, accumulate, parent_id, company_id from vw_budget_item where id = $1 and company_id = $2"
	err := s.db.QueryRow(sql, id, companyId).Scan(&bi.ID, &bi.Code, &bi.Name, &bi.Level, &bi.Accumulate, &bi.ParentId, &bi.CompanyId)

	return bi, err
}

func (s *service) UpdateBudgetItem(bi *types.BudgetItem) error {
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

	sql := "update budget_item set code = $1, name = $2, level = $3, accumulate = $4, parent_id = $5 where id = $6 and company_id = $7"
	_, err := s.db.Exec(sql, bi.Code, bi.Name, level, bi.Accumulate, bi.ParentId, bi.ID, bi.CompanyId)
	return err
}

func (s *service) getBudgetItem(id, companyId uuid.UUID) (*types.BudgetItem, error) {
	if id == uuid.Nil {
		return nil, nil
	}

	sql := `
		  select id, code, name, level, accumulate, parent_id, company_id
		  from vw_budget_item where id = $1 and company_id = $2
	 `

	bi := &types.BudgetItem{}
	err := s.db.QueryRow(
		sql,
		id,
		companyId,
	).Scan(
		&bi.ID,
		&bi.Code,
		&bi.Name,
		&bi.Level,
		&bi.Accumulate,
		&bi.ParentId,
		&bi.CompanyId,
	)

	return bi, err
}

func (s *service) GetBudgetItemsByAccumulate(companyId uuid.UUID, accum bool) []types.BudgetItem {
	sql := `
		  select id, code, name, level, accumulate, parent_id, company_id
		  from vw_budget_item
		  where company_id = $1 and accumulate = $2 order by name
		  `

	rows, err := s.db.Query(sql, companyId, accum)
	if err != nil {
		return nil
	}
	defer rows.Close()
	bis := []types.BudgetItem{}

	for rows.Next() {
		bi := types.BudgetItem{}
		if err := rows.Scan(&bi.ID, &bi.Code, &bi.Name, &bi.Level, &bi.Accumulate, &bi.ParentId, &bi.CompanyId); err != nil {
			return nil
		}
		bis = append(bis, bi)
	}

	return bis
}
