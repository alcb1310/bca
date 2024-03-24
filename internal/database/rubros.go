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

func (s *service) CreateRubro(rubro types.Rubro) error {
	query := "insert into item (code, name, unit, company_id) values ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, rubro.Code, rubro.Name, rubro.Unit, rubro.CompanyId)
	return err
}
