package database

import (
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
)

func (s *service) GetAllCategories(companyId uuid.UUID) ([]types.Category, error) {
	categories := []types.Category{}
	sql := "select id, name from category where company_id = $1 order by name"
	rows, err := s.db.Query(sql, companyId)
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		c := types.Category{}
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			return categories, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (s *service) CreateCategory(category types.Category) error {
	sql := "insert into category (name, company_id) values ($1, $2)"
	_, err := s.db.Exec(sql, category.Name, category.CompanyId)
	return err
}

func (s *service) GetCategory(id, companyId uuid.UUID) (types.Category, error) {
	c := types.Category{}
	sql := "select id, name, company_id from category where id = $1 and company_id = $2"
	err := s.db.QueryRow(sql, id, companyId).Scan(&c.Id, &c.Name, &c.CompanyId)
	return c, err
}

func (s *service) UpdateCategory(category types.Category) error {
	sql := "update category set name = $1 where id = $2 and company_id = $3"
	_, err := s.db.Exec(sql, category.Name, category.Id, category.CompanyId)
	return err
}
