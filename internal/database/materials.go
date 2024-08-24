package database

import (
	"log"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
)

func (s *service) GetAllMaterials(companyId uuid.UUID) []types.Material {
	materials := []types.Material{}
	query := `
    select
      id, code, name, unit, category_name, category_id, company_id
    from
      vw_materials
    where company_id = $1
    order by category_name, code
    `
	rows, err := s.db.Query(query, companyId)
	if err != nil {
		log.Fatal(err)
		return materials
	}
	defer rows.Close()

	for rows.Next() {
		var material types.Material
		if err := rows.Scan(&material.Id, &material.Code, &material.Name, &material.Unit, &material.Category.Name, &material.Category.Id, &material.CompanyId); err != nil {
			log.Fatal(err)
			return materials
		}
		materials = append(materials, material)
	}

	return materials
}

func (s *service) CreateMaterial(material types.Material) error {
	sql := "insert into materials (code, name, unit, category_id, company_id) values ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(sql, material.Code, material.Name, material.Unit, material.Category.Id, material.CompanyId)
	return err
}

func (s *service) GetMaterial(id, companyId uuid.UUID) (types.Material, error) {
	material := types.Material{}
	query := `
    select
      id, code, name, unit, category_name, category_id, company_id
    from
      vw_materials
    where id = $1 and company_id = $2
    `
	err := s.db.QueryRow(query, id, companyId).Scan(&material.Id, &material.Code, &material.Name, &material.Unit, &material.Category.Name, &material.Category.Id, &material.CompanyId)
	return material, err
}

func (s *service) UpdateMaterial(material types.Material) error {
	sql := "update materials set code = $1, name = $2, unit = $3, category_id = $4 where id = $5 and company_id = $6"
	_, err := s.db.Exec(sql, material.Code, material.Name, material.Unit, material.Category.Id, material.Id, material.CompanyId)
	return err
}
