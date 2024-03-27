package database

import (
	"log"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *service) CreateCantidades(projectId, rubroId uuid.UUID, quantity float64, companyId uuid.UUID) error {
	query := "insert into analysis (project_id, item_id, quantity, company_id) values ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, projectId, rubroId, quantity, companyId)
	return err
}

func (s *service) CantidadesTable(companyId uuid.UUID) []types.Quantity {
	quantities := []types.Quantity{}

	query := `
    select
      id, quantity, project_id, project_name,
      item_id, item_code, item_name, item_unit
    from vw_project_costs
    where company_id = $1
    order by project_name, item_name
  `
	rows, err := s.db.Query(query, companyId)
	if err != nil {
		log.Println(err)
		return quantities
	}
	defer rows.Close()

	for rows.Next() {
		var quantity types.Quantity
		if err := rows.Scan(&quantity.Id, &quantity.Quantity, &quantity.Project.ID, &quantity.Project.Name,
			&quantity.Rubro.Id, &quantity.Rubro.Code, &quantity.Rubro.Name, &quantity.Rubro.Unit); err != nil {
			log.Fatal(err)
			return quantities
		}
		quantities = append(quantities, quantity)
	}

	return quantities
}

func (s *service) AnalysisReport(project_id, company_id uuid.UUID) []types.AnalysisReport {
	analysis := []types.AnalysisReport{}

	sql := `
    select project_name, category_name, material_name, sum(quantity * item_material_quantity)
    from vw_project_cost_analysis
    where project_id = $1 and company_id = $2
    group by project_name, category_name, material_name
    order by project_name, category_name, material_name
  `

	rows, err := s.db.Query(sql, project_id, company_id)
	if err != nil {
		log.Println(err)
		return analysis
	}
	defer rows.Close()

	for rows.Next() {
		var ar types.AnalysisReport
		if err := rows.Scan(&ar.ProjectName, &ar.CategoryName, &ar.MaterialName, &ar.Quantity); err != nil {
			log.Fatal(err)
			return analysis
		}
		analysis = append(analysis, ar)
	}

	return analysis
}
