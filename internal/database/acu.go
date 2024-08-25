package database

import (
	"log/slog"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
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
		slog.Error(err.Error())
		return quantities
	}
	defer rows.Close()

	for rows.Next() {
		var quantity types.Quantity
		if err := rows.Scan(&quantity.Id, &quantity.Quantity, &quantity.Project.ID, &quantity.Project.Name,
			&quantity.Rubro.Id, &quantity.Rubro.Code, &quantity.Rubro.Name, &quantity.Rubro.Unit); err != nil {
			slog.Error(err.Error())
			return quantities
		}
		quantities = append(quantities, quantity)
	}

	return quantities
}

func (s *service) AnalysisReport(project_id, company_id uuid.UUID) map[string][]types.AnalysisReport {
	x := make(map[string][]types.AnalysisReport)

	sql := `
    select project_name, category_name, material_name, sum(quantity * item_material_quantity)
    from vw_project_cost_analysis
    where project_id = $1 and company_id = $2
    group by project_name, category_name, material_name
    order by project_name, category_name, material_name
  `

	rows, err := s.db.Query(sql, project_id, company_id)
	if err != nil {
		slog.Error(err.Error())
		return x
	}
	defer rows.Close()

	for rows.Next() {
		var analysis types.AnalysisReport
		if err := rows.Scan(&analysis.ProjectName, &analysis.CategoryName, &analysis.MaterialName, &analysis.Quantity); err != nil {
			slog.Error(err.Error())
			return x
		}

		_, ok := x[analysis.CategoryName]
		if ok {
			x[analysis.CategoryName] = append(x[analysis.CategoryName], analysis)
		} else {
			x[analysis.CategoryName] = []types.AnalysisReport{analysis}
		}
	}

	return x
}

func (s *service) GetQuantityByMaterialAndItem(itemId, materialId, companyId uuid.UUID) types.ItemMaterialType {
	itemMaterial := types.ItemMaterialType{}
	query := "select quantity from item_materials where item_id = $1 and material_id = $2 and company_id = $3"
	err := s.db.QueryRow(query, itemId, materialId, companyId).Scan(&itemMaterial.Quantity)
	if err != nil {
		slog.Error(err.Error())
		return itemMaterial
	}

	itemMaterial.ItemId = itemId
	itemMaterial.MaterialId = materialId

	return itemMaterial
}

func (s *service) DeleteCantidades(id, companyId uuid.UUID) error {
	query := "delete from analysis where id = $1 and company_id = $2"
	_, err := s.db.Exec(query, id, companyId)
	return err
}

func (s *service) GetOneQuantityById(id, companyId uuid.UUID) types.Quantity {
	quantity := types.Quantity{}
	query := "select id, quantity, project_id, project_name, item_id, item_code, item_name, item_unit from vw_project_costs where id = $1 and company_id = $2"
	err := s.db.QueryRow(query, id, companyId).Scan(&quantity.Id, &quantity.Quantity, &quantity.Project.ID, &quantity.Project.Name,
		&quantity.Rubro.Id, &quantity.Rubro.Code, &quantity.Rubro.Name, &quantity.Rubro.Unit)
	if err != nil {
		slog.Error(err.Error())
		return quantity
	}
	return quantity
}

func (s *service) UpdateQuantity(q types.Quantity, companyId uuid.UUID) error {
	query := "update analysis set quantity = $1 where id = $2 and company_id = $3"
	_, err := s.db.Exec(query, q.Quantity, q.Id, companyId)
	return err
}
