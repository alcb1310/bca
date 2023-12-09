package database

import (
	"bca-go-final/internal/types"
	"log"

	"github.com/google/uuid"
)

func (s *service) GetAllUsers(companyId uuid.UUID) ([]types.User, error) {
	users := []types.User{}

	sql := "select id, name, email, company_id, role_id from \"user\" where company_id = $1"
	rows, err := s.db.Query(sql, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := types.User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.CompanyId, &u.RoleId); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	log.Println("GetAllUsers: ", users)

	return users, nil
}
