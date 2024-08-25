package database

import (
	"log/slog"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
)

func (s *service) GetAllRubros(companyId uuid.UUID) ([]types.Rubro, error) {
	rubros := []types.Rubro{}
	query := "select id, code, name, unit, company_id from item where company_id = $1 order by name"
	rows, err := s.db.Query(query, companyId)
	if err != nil {
		slog.Error("Error getting rubros: ", "err", err)
		return rubros, err
	}
	defer rows.Close()

	for rows.Next() {
		var rubro types.Rubro
		if err := rows.Scan(&rubro.Id, &rubro.Code, &rubro.Name, &rubro.Unit, &rubro.CompanyId); err != nil {
			slog.Error("Error getting rubros: ", "err", err)
			return rubros, err
		}
		rubros = append(rubros, rubro)
	}

	return rubros, nil
}

func (s *service) CreateRubro(rubro types.Rubro) (uuid.UUID, error) {
	var id uuid.UUID
	query := "insert into item (code, name, unit, company_id) values ($1, $2, $3, $4) returning id"
	err := s.db.QueryRow(query, &rubro.Code, &rubro.Name, &rubro.Unit, &rubro.CompanyId).Scan(&id)
	return id, err
}

func (s *service) GetOneRubro(id, companyId uuid.UUID) (types.Rubro, error) {
	rubro := types.Rubro{}
	query := "select id, code, name, unit, company_id from item where id = $1 and company_id = $2"
	err := s.db.QueryRow(query, id, companyId).Scan(&rubro.Id, &rubro.Code, &rubro.Name, &rubro.Unit, &rubro.CompanyId)
	return rubro, err
}

func (s *service) UpdateRubro(rubro types.Rubro) error {
	query := "update item set code = $1, name = $2, unit = $3 where id = $4 and company_id = $5"
	_, err := s.db.Exec(query, rubro.Code, rubro.Name, rubro.Unit, rubro.Id, rubro.CompanyId)
	return err
}

func (s *service) GetMaterialsByItem(id, companyId uuid.UUID) []types.ACU {
	acus := []types.ACU{}

	query := `
    select
      item_id, item_code, item_name, item_unit,
      material_id, material_code, material_name, material_unit,
      quantity, company_id
    from vw_acu
    where
      item_id = $1 and company_id = $2
    `

	rows, err := s.db.Query(query, id, companyId)
	if err != nil {
		slog.Error("Error getting materials by item: ", "err", err)
		return acus
	}
	defer rows.Close()

	for rows.Next() {
		var acu types.ACU
		if err := rows.Scan(&acu.Item.Id, &acu.Item.Code, &acu.Item.Name, &acu.Item.Unit, &acu.Material.Id, &acu.Material.Code, &acu.Material.Name, &acu.Material.Unit, &acu.Quantity, &acu.CompanyId); err != nil {
			slog.Error("Error getting materials by item: ", "err", err)
			return acus
		}
		acus = append(acus, acu)
	}

	return acus
}

func (s *service) AddMaterialsByItem(itemId, materialId uuid.UUID, quantity float64, companyId uuid.UUID) error {
	query := "insert into item_materials (item_id, material_id, quantity, company_id) values ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, itemId, materialId, quantity, companyId)
	return err
}

func (s *service) DeleteMaterialsByItem(itemId, materialId, companyId uuid.UUID) error {
	query := "delete from item_materials where item_id = $1 and material_id = $2 and company_id = $3"
	_, err := s.db.Exec(query, itemId, materialId, companyId)
	return err
}

func (s *service) UpdateMaterialByItem(itemId, materialId uuid.UUID, quantity float64, companyId uuid.UUID) error {
	query := "update item_materials set quantity = $1 where item_id = $2 and material_id = $3 and company_id = $4"
	_, err := s.db.Exec(query, quantity, itemId, materialId, companyId)
	return err
}
