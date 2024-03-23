package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetAllCategories(companyId uuid.UUID) ([]types.Category, error) {
  categories := []types.Category{}
  sql := "select id, name from category where company_id = $1"
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
