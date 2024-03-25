package database

import (
	"log"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *service) GetAllRubros(companyId uuid.UUID) ([]types.Rubro, error) {

	rubros := []types.Rubro{}
	query := "select id, code, name, unit, company_id from item where company_id = $1 order by name"
	rows, err := s.db.Query(query, companyId)
	if err != nil {
		log.Println("Error getting rubros: ", err)
		return rubros, err
	}
	defer rows.Close()

	for rows.Next() {
		var rubro types.Rubro
		if err := rows.Scan(&rubro.Id, &rubro.Code, &rubro.Name, &rubro.Unit, &rubro.CompanyId); err != nil {
			log.Println("Error getting rubros: ", err)
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
