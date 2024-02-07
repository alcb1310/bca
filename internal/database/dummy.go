package database

import (
	"database/sql"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

var (
	accumTrue  = true
	accumFalse = false
	cdcId      = uuid.Nil
	ggId       = uuid.Nil
	ogId       = uuid.Nil
	haId       = uuid.Nil
	chId       = uuid.Nil
)

/*
func (s *service) LoadDummyData(companyId uuid.UUID) error {
	log.Println("Loading dummy data")

	var count int
	query := "select count(*) from budget_item where company_id = $1"
	err := s.db.QueryRow(query, companyId).Scan(&count)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}
	if count > 0 {
		log.Println("Dummy data already loaded")
		return nil
	}

	tx, _ := s.db.Begin()
	defer tx.Rollback()
	bi := &types.BudgetItem{
		CompanyId:  companyId,
		Code:       "500",
		Name:       "COSTO DIRECTO DE CONSTRUCCION",
		Accumulate: &accumTrue,
		ParentId:   nil,
	}

	cdcId, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	bi.Code = "500.01"
	bi.Name = "GASTOS GENERALES"
	bi.ParentId = &cdcId
	ggId, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	bi.Code = "500.01.01"
	bi.Name = "BODEGUERO"
	bi.ParentId = &ggId
	bi.Accumulate = &accumFalse
	_, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	bi.Code = "500.01.02"
	bi.Name = "WACHIMAN"
	bi.ParentId = &ggId
	bi.Accumulate = &accumFalse
	_, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	bi.Code = "500.02"
	bi.Name = "OBRA GRUESA"
	bi.ParentId = &cdcId
	bi.Accumulate = &accumTrue
	ogId, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	bi.Code = "500.02.01"
	bi.Name = "HIERRO Y ALAMBRE"
	bi.ParentId = &ogId
	bi.Accumulate = &accumTrue
	haId, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	bi.Code = "500.02.01.01"
	bi.Name = "VARILLA 08MM X 12M"
	bi.ParentId = &haId
	bi.Accumulate = &accumFalse
	_, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	bi.Code = "500.02.01.02"
	bi.Name = "VARILLA 10MM X 12M"
	bi.ParentId = &haId
	bi.Accumulate = &accumFalse
	_, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	bi.Code = "500.02.02"
	bi.Name = "CEMENTO Y HORMIGONES"
	bi.ParentId = &ogId
	bi.Accumulate = &accumTrue
	haId, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	bi.Code = "500.02.02.01"
	bi.Name = "CEMENTO"
	bi.ParentId = &haId
	bi.Accumulate = &accumFalse
	_, err = createBudgetItem(tx, bi)
	if err != nil {
		log.Println("Error loading dummy data")
		return err
	}

	tx.Commit()
	log.Println("Dummy data loaded")

	return nil
} */

func createBudgetItem(tx *sql.Tx, bi *types.BudgetItem) (uuid.UUID, error) {
	createdId := uuid.Nil
	var level uint8 = 1
	if bi.ParentId != nil {
		sql := "select level from budget_item where id = $1 and company_id = $2"
		err := tx.QueryRow(sql, bi.ParentId, bi.CompanyId).Scan(&level)
		if err != nil {
			return createdId, err
		}
		level++
	}
	bi.Level = level

	sql := "insert into budget_item (code, name, level, accumulate, parent_id, company_id) values ($1, $2, $3, $4, $5, $6) returning id"
	err := tx.QueryRow(sql, bi.Code, bi.Name, bi.Level, bi.Accumulate, bi.ParentId, bi.CompanyId).Scan(&createdId)

	return createdId, err
}
